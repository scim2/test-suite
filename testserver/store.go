package testserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"math/rand/v2"
	"net/http"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

// writeOnlyAttrs lists attribute names that must never be returned.
var writeOnlyAttrs = []string{"password", "writeOnlyAttr"}

// enrichMembers adds $ref URIs to group member sub-attributes.
func enrichMembers(attrs scim.ResourceAttributes) {
	members, ok := attrs["members"].([]any)
	if !ok {
		return
	}
	for i, m := range members {
		mm, ok := m.(map[string]any)
		if !ok {
			continue
		}
		if _, hasRef := mm["$ref"]; hasRef {
			continue
		}
		val, _ := mm["value"].(string)
		typ, _ := mm["type"].(string)
		if val == "" {
			continue
		}
		endpoint := "Users"
		if typ == "Group" {
			endpoint = "Groups"
		}
		mm["$ref"] = fmt.Sprintf("/%s/%s", endpoint, val)
		members[i] = mm
	}
	attrs["members"] = members
}

func externalID(attrs scim.ResourceAttributes) optional.String {
	if v, ok := attrs["externalId"].(string); ok {
		return optional.NewString(v)
	}
	return optional.String{}
}

func newID() string {
	return fmt.Sprintf("%08x%08x", rand.Uint32(), rand.Uint32())
}

func scanResource(rows *sql.Rows) (scim.Resource, error) {
	var (
		id                string
		rawAttrs          string
		created, modified time.Time
		version           string
	)
	if err := rows.Scan(&id, &rawAttrs, &created, &modified, &version); err != nil {
		return scim.Resource{}, err
	}

	var attrs scim.ResourceAttributes
	if err := json.Unmarshal([]byte(rawAttrs), &attrs); err != nil {
		return scim.Resource{}, err
	}

	return toResource(id, attrs, created, modified, version), nil
}

func toResource(id string, attrs scim.ResourceAttributes, created, modified time.Time, version string) scim.Resource {
	// Strip writeOnly attributes (e.g. password, writeOnlyAttr).
	clean := make(scim.ResourceAttributes, len(attrs))
	maps.Copy(clean, attrs)
	for _, name := range writeOnlyAttrs {
		delete(clean, name)
	}

	// Enrich group members with $ref.
	enrichMembers(clean)

	return scim.Resource{
		ID:         id,
		ExternalID: externalID(attrs),
		Attributes: clean,
		Meta: scim.Meta{
			Created:      &created,
			LastModified: &modified,
			Version:      version,
		},
	}
}

func uniqueName(attrs scim.ResourceAttributes) string {
	if v, ok := attrs["userName"].(string); ok {
		return v
	}
	if v, ok := attrs["displayName"].(string); ok {
		return v
	}
	return ""
}

// handler implements scim.ResourceHandler backed by the shared store.
type handler struct {
	store      *store
	tableName  string
	schema     schema.Schema
	extensions []schema.Schema
}

func newHandler(s *store, tableName string, sc schema.Schema, extensions ...schema.Schema) *handler {
	return &handler{store: s, tableName: tableName, schema: sc, extensions: extensions}
}

func (h *handler) Create(_ *http.Request, attrs scim.ResourceAttributes) (scim.Resource, error) {
	// Check uniqueness.
	name := uniqueName(attrs)
	if name != "" {
		var count int
		err := h.store.db.QueryRow(
			`SELECT count(*) FROM resources
			 WHERE table_name = $1
			   AND json_extract_string(attrs, '$.userName') = $2
			    OR json_extract_string(attrs, '$.displayName') = $2`,
			h.tableName, name,
		).Scan(&count)
		if err != nil {
			return scim.Resource{}, errors.ScimErrorInternal
		}
		if count > 0 {
			return scim.Resource{}, errors.ScimErrorUniqueness
		}
	}

	id := newID()
	now := time.Now().UTC()
	ver := fmt.Sprintf("W/\"%d\"", now.UnixNano())

	data, err := json.Marshal(attrs)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	_, err = h.store.db.Exec(
		`INSERT INTO resources (id, table_name, attrs, created, modified, version)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, h.tableName, string(data), now, now, ver,
	)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	return toResource(id, attrs, now, now, ver), nil
}

func (h *handler) Delete(_ *http.Request, id string) error {
	res, err := h.store.db.Exec(
		`DELETE FROM resources WHERE table_name = $1 AND id = $2`,
		h.tableName, id,
	)
	if err != nil {
		return errors.ScimErrorInternal
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.ScimErrorResourceNotFound(id)
	}
	return nil
}

func (h *handler) Get(_ *http.Request, id string) (scim.Resource, error) {
	return h.getByID(id)
}

func (h *handler) GetAll(_ *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	rows, err := h.store.db.Query(
		`SELECT id, attrs, created, modified, version
		 FROM resources WHERE table_name = $1
		 ORDER BY created`,
		h.tableName,
	)
	if err != nil {
		return scim.Page{}, errors.ScimErrorInternal
	}
	defer func() { _ = rows.Close() }()

	var all []scim.Resource
	for rows.Next() {
		r, err := scanResource(rows)
		if err != nil {
			return scim.Page{}, errors.ScimErrorInternal
		}
		if params.FilterValidator != nil {
			if err := params.FilterValidator.PassesFilter(r.Attributes); err != nil {
				continue
			}
		}
		all = append(all, r)
	}

	total := len(all)

	start := min(max(params.StartIndex-1, 0), total)

	end := total
	if params.Count > 0 && start+params.Count < end {
		end = start + params.Count
	}

	return scim.Page{
		TotalResults: total,
		Resources:    all[start:end],
	}, nil
}

func (h *handler) Patch(_ *http.Request, id string, operations []scim.PatchOperation) (scim.Resource, error) {
	existing, err := h.getByID(id)
	if err != nil {
		return scim.Resource{}, err
	}

	attrs, err := scim.ApplyPatch(existing.Attributes, operations, h.schema, h.extensions...)
	if err != nil {
		return scim.Resource{}, err
	}

	now := time.Now().UTC()
	ver := fmt.Sprintf("W/\"%d\"", now.UnixNano())

	data, err := json.Marshal(attrs)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	_, err = h.store.db.Exec(
		`UPDATE resources SET attrs = $1, modified = $2, version = $3
		 WHERE table_name = $4 AND id = $5`,
		string(data), now, ver, h.tableName, id,
	)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	created := existing.Meta.Created
	return toResource(id, attrs, *created, now, ver), nil
}

func (h *handler) Replace(r *http.Request, id string, attrs scim.ResourceAttributes) (scim.Resource, error) {
	// Check exists.
	existing, err := h.getByID(id)
	if err != nil {
		return scim.Resource{}, err
	}

	// Check ETag precondition.
	if ifMatch := r.Header.Get("If-Match"); ifMatch != "" && ifMatch != existing.Meta.Version {
		return scim.Resource{}, errors.ScimError{Status: http.StatusPreconditionFailed}
	}

	// Check uniqueness if name changed.
	newName := uniqueName(attrs)
	oldName := uniqueName(existing.Attributes)
	if newName != "" && newName != oldName {
		var count int
		err := h.store.db.QueryRow(
			`SELECT count(*) FROM resources
			 WHERE table_name = $1 AND id != $2
			   AND (json_extract_string(attrs, '$.userName') = $3
			    OR  json_extract_string(attrs, '$.displayName') = $3)`,
			h.tableName, id, newName,
		).Scan(&count)
		if err != nil {
			return scim.Resource{}, errors.ScimErrorInternal
		}
		if count > 0 {
			return scim.Resource{}, errors.ScimErrorUniqueness
		}
	}

	now := time.Now().UTC()
	ver := fmt.Sprintf("W/\"%d\"", now.UnixNano())

	data, err := json.Marshal(attrs)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	_, err = h.store.db.Exec(
		`UPDATE resources SET attrs = $1, modified = $2, version = $3
		 WHERE table_name = $4 AND id = $5`,
		string(data), now, ver, h.tableName, id,
	)
	if err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	created := existing.Meta.Created
	return toResource(id, attrs, *created, now, ver), nil
}

func (h *handler) Search(r *http.Request, params scim.SearchParams) (scim.Page, error) {
	return h.GetAll(r, scim.ListRequestParams{
		Count:      params.Count,
		StartIndex: params.StartIndex,
	})
}

func (h *handler) getByID(id string) (scim.Resource, error) {
	row := h.store.db.QueryRow(
		`SELECT id, attrs, created, modified, version
		 FROM resources WHERE table_name = $1 AND id = $2`,
		h.tableName, id,
	)

	var (
		rid               string
		rawAttrs          string
		created, modified time.Time
		version           string
	)
	if err := row.Scan(&rid, &rawAttrs, &created, &modified, &version); err != nil {
		return scim.Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	var attrs scim.ResourceAttributes
	if err := json.Unmarshal([]byte(rawAttrs), &attrs); err != nil {
		return scim.Resource{}, errors.ScimErrorInternal
	}

	return toResource(rid, attrs, created, modified, version), nil
}

// rootQueryHandler implements scim.RootQueryHandler and
// scim.ResourceSearcher by querying all resource types in the store.
type rootQueryHandler struct {
	store *store
}

func (h *rootQueryHandler) GetAll(_ *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	return h.query(params.Count, params.StartIndex)
}

func (h *rootQueryHandler) Search(_ *http.Request, params scim.SearchParams) (scim.Page, error) {
	return h.query(params.Count, params.StartIndex)
}

func (h *rootQueryHandler) query(count, startIndex int) (scim.Page, error) {
	rows, err := h.store.db.Query(
		`SELECT id, attrs, created, modified, version
		 FROM resources ORDER BY created`,
	)
	if err != nil {
		return scim.Page{}, errors.ScimErrorInternal
	}
	defer func() { _ = rows.Close() }()

	var all []scim.Resource
	for rows.Next() {
		r, err := scanResource(rows)
		if err != nil {
			return scim.Page{}, errors.ScimErrorInternal
		}
		all = append(all, r)
	}

	total := len(all)
	start := min(max(startIndex-1, 0), total)
	end := total
	if count > 0 && start+count < end {
		end = start + count
	}

	return scim.Page{
		TotalResults: total,
		Resources:    all[start:end],
	}, nil
}

// store is a DuckDB-backed resource store shared across handlers.
type store struct {
	db *sql.DB
}

func newStore() (*store, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("open duckdb: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE resources (
			id         TEXT NOT NULL,
			table_name TEXT NOT NULL,
			attrs      TEXT NOT NULL,
			created    TIMESTAMP NOT NULL,
			modified   TIMESTAMP NOT NULL,
			version    TEXT NOT NULL,
			PRIMARY KEY (table_name, id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &store{db: db}, nil
}
