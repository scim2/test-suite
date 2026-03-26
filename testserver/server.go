// Package testserver provides an in-memory SCIM 2.0 server for testing.
// It uses github.com/elimity-com/scim for the protocol layer and DuckDB
// as an in-memory store for resources. The server wraps the library's
// handler to add /Bulk, /Me, and sorting support.
package testserver

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

const (
	defaultBulkMaxOps     = 10
	defaultBulkMaxPayload = 1048576 // 1 MB
)

// New creates an http.Handler serving a SCIM 2.0 API backed by DuckDB.
func New() http.Handler {
	store, err := newStore()
	if err != nil {
		panic(fmt.Sprintf("testserver: %v", err))
	}

	config := scim.ServiceProviderConfig{
		SupportFiltering: true,
		SupportPatch:     true,
		MaxResults:       50,
		AuthenticationSchemes: []scim.AuthenticationScheme{
			{
				Type:        scim.AuthenticationTypeOauthBearerToken,
				Name:        "OAuth Bearer Token",
				Description: "Authentication scheme using the OAuth Bearer Token Standard",
			},
		},
	}

	resourceTypes := []scim.ResourceType{
		{
			ID:          optional.NewString("User"),
			Name:        "User",
			Endpoint:    "/Users",
			Description: optional.NewString("User Account"),
			Schema:      schema.CoreUserSchema(),
			Handler:     newHandler(store, "users", schema.CoreUserSchema()),
		},
		{
			ID:          optional.NewString("Group"),
			Name:        "Group",
			Endpoint:    "/Groups",
			Description: optional.NewString("Group"),
			Schema:      schema.CoreGroupSchema(),
			Handler:     newHandler(store, "groups", schema.CoreGroupSchema()),
		},
		{
			ID:          optional.NewString("TestResource"),
			Name:        "TestResource",
			Endpoint:    "/TestResources",
			Description: optional.NewString("Synthetic resource for compliance testing"),
			Schema:      testResourceSchema(),
			SchemaExtensions: []scim.SchemaExtension{
				{
					Schema:   testExtensionSchema(),
					Required: true,
				},
			},
			Handler: newHandler(store, "testresources", testResourceSchema(), testExtensionSchema()),
		},
	}

	server, err := scim.NewServer(
		&scim.ServerArgs{
			ServiceProviderConfig: &config,
			ResourceTypes:         resourceTypes,
		},
		scim.WithRootQueryHandler(&rootQueryHandler{store: store}),
	)
	if err != nil {
		panic(fmt.Sprintf("testserver: %v", err))
	}

	return &wrapper{
		inner:          server,
		store:          store,
		bulkMaxOps:     defaultBulkMaxOps,
		bulkMaxPayload: defaultBulkMaxPayload,
	}
}

// attrString extracts a string attribute from a resource map,
// supporting dotted sub-attribute paths like "name.givenName".
func attrString(m map[string]any, attr string) string {
	if m == nil {
		return ""
	}
	parts := strings.SplitN(attr, ".", 2)
	if len(parts) == 2 {
		sub, ok := m[parts[0]].(map[string]any)
		if !ok {
			return ""
		}
		v, _ := sub[parts[1]].(string)
		return v
	}
	v, _ := m[attr].(string)
	return v
}

// copyResponse writes a recorded response to the real writer.
func copyResponse(rw http.ResponseWriter, rec *httptest.ResponseRecorder) {
	maps.Copy(rw.Header(), rec.Header())
	rw.WriteHeader(rec.Code)
	_, _ = rw.Write(rec.Body.Bytes())
}

// writeSCIMError writes a SCIM error response.
func writeSCIMError(w http.ResponseWriter, status int, scimType string) {
	resp := map[string]any{
		"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:Error"},
		"status":  fmt.Sprintf("%d", status),
	}
	if scimType != "" {
		resp["scimType"] = scimType
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// wrapper intercepts requests to add /Bulk, /Me, and sorting support
// on top of the elimity-com/scim library handler.
type wrapper struct {
	inner          http.Handler
	store          *store
	bulkMaxOps     int
	bulkMaxPayload int
}

func (w *wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/ServiceProviderConfig" && r.Method == http.MethodGet:
		w.handleServiceProviderConfig(rw, r)
		return
	case path == "/Bulk" && r.Method == http.MethodPost:
		w.handleBulk(rw, r)
		return
	case path == "/Me" && r.Method == http.MethodGet:
		w.handleMe(rw, r)
		return
	}

	if r.Method == http.MethodGet && r.URL.Query().Get("sortBy") != "" {
		w.handleSortedList(rw, r)
		return
	}

	w.inner.ServeHTTP(rw, r)
}

// handleMe returns the first user as the "current" user, with
// Location and Content-Location headers pointing to the permanent URI.
func (w *wrapper) handleMe(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/scim+json")

	row := w.store.db.QueryRow(
		`SELECT id FROM resources WHERE table_name = 'users' ORDER BY created LIMIT 1`,
	)
	var id string
	if err := row.Scan(&id); err != nil {
		writeSCIMError(rw, 404, "")
		return
	}

	subReq, _ := http.NewRequest("GET", "/Users/"+id, nil)
	rec := httptest.NewRecorder()
	w.inner.ServeHTTP(rec, subReq)

	maps.Copy(rw.Header(), rec.Header())
	location := "/Users/" + id
	rw.Header().Set("Location", location)
	rw.Header().Set("Content-Location", location)
	rw.WriteHeader(rec.Code)
	_, _ = rw.Write(rec.Body.Bytes())
}

// handleServiceProviderConfig intercepts the config response to
// advertise bulk and sort support.
func (w *wrapper) handleServiceProviderConfig(rw http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	w.inner.ServeHTTP(rec, r)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		copyResponse(rw, rec)
		return
	}

	body["bulk"] = map[string]any{
		"supported":      true,
		"maxOperations":  w.bulkMaxOps,
		"maxPayloadSize": w.bulkMaxPayload,
	}
	body["sort"] = map[string]any{
		"supported": true,
	}
	body["etag"] = map[string]any{
		"supported": true,
	}
	body["changePassword"] = map[string]any{
		"supported": true,
	}

	data, _ := json.Marshal(body)
	rw.Header().Set("Content-Type", "application/scim+json")
	_, _ = rw.Write(data)
}

// handleSortedList captures a list response and sorts the Resources
// array by the requested sortBy attribute.
func (w *wrapper) handleSortedList(rw http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	w.inner.ServeHTTP(rec, r)

	if rec.Code != 200 {
		copyResponse(rw, rec)
		return
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		copyResponse(rw, rec)
		return
	}

	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")
	if sortOrder == "" {
		sortOrder = "ascending"
	}

	resources, ok := body["Resources"].([]any)
	if ok && len(resources) > 1 {
		slices.SortStableFunc(resources, func(a, b any) int {
			am, _ := a.(map[string]any)
			bm, _ := b.(map[string]any)
			av := attrString(am, sortBy)
			bv := attrString(bm, sortBy)
			cmp := strings.Compare(strings.ToLower(av), strings.ToLower(bv))
			if sortOrder == "descending" {
				cmp = -cmp
			}
			return cmp
		})
		body["Resources"] = resources
	}

	data, _ := json.Marshal(body)
	maps.Copy(rw.Header(), rec.Header())
	rw.WriteHeader(200)
	_, _ = rw.Write(data)
}
