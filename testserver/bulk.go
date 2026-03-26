package testserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	bulkResponseSchema = "urn:ietf:params:scim:api:messages:2.0:BulkResponse"
)

func parseStatusCode(s string) int {
	code := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			code = code*10 + int(c-'0')
		}
	}
	return code
}

// resolveBulkIDs recursively replaces "bulkId:xxx" references with
// the actual resource IDs of previously created resources.
func resolveBulkIDs(data map[string]any, bulkIDs map[string]string) map[string]any {
	if data == nil {
		return nil
	}
	resolved := make(map[string]any, len(data))
	for k, v := range data {
		resolved[k] = resolveValue(v, bulkIDs)
	}
	return resolved
}

func resolveValue(v any, bulkIDs map[string]string) any {
	switch val := v.(type) {
	case string:
		if after, ok := strings.CutPrefix(val, "bulkId:"); ok {
			ref := after
			if id, ok := bulkIDs[ref]; ok {
				return id
			}
		}
		return val
	case map[string]any:
		resolved := make(map[string]any, len(val))
		for k, item := range val {
			resolved[k] = resolveValue(item, bulkIDs)
		}
		return resolved
	case []any:
		resolved := make([]any, len(val))
		for i, item := range val {
			resolved[i] = resolveValue(item, bulkIDs)
		}
		return resolved
	default:
		return v
	}
}

type bulkOp struct {
	Method string         `json:"method"`
	Path   string         `json:"path"`
	BulkID string         `json:"bulkId"`
	Data   map[string]any `json:"data"`
}

type bulkRequest struct {
	Schemas      []string `json:"schemas"`
	FailOnErrors int      `json:"failOnErrors"`
	Operations   []bulkOp `json:"Operations"`
}

func (w *wrapper) executeBulkOp(op bulkOp, data map[string]any) map[string]any {
	result := map[string]any{
		"method": op.Method,
	}

	method := strings.ToUpper(op.Method)

	var subReq *http.Request
	var err error

	switch method {
	case "POST", "PUT", "PATCH":
		encoded, _ := json.Marshal(data)
		subReq, err = http.NewRequest(method, op.Path, bytes.NewReader(encoded))
		if subReq != nil {
			subReq.Header.Set("Content-Type", "application/scim+json")
		}
	case "DELETE":
		subReq, err = http.NewRequest("DELETE", op.Path, nil)
	default:
		result["status"] = "400"
		return result
	}

	if err != nil {
		result["status"] = "500"
		return result
	}

	rec := httptest.NewRecorder()
	w.inner.ServeHTTP(rec, subReq)

	result["status"] = fmt.Sprintf("%d", rec.Code)

	// Extract location for created resources.
	if method == "POST" && rec.Code == 201 {
		var respBody map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err == nil {
			if id, _ := respBody["id"].(string); id != "" {
				result["location"] = op.Path + "/" + id
			}
		}
	}

	return result
}

func (w *wrapper) handleBulk(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/scim+json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeSCIMError(rw, 400, "invalidSyntax")
		return
	}

	var req bulkRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeSCIMError(rw, 400, "invalidSyntax")
		return
	}

	if len(req.Operations) > w.bulkMaxOps {
		writeSCIMError(rw, 413, "")
		return
	}

	if len(body) > w.bulkMaxPayload {
		writeSCIMError(rw, 413, "")
		return
	}

	bulkIDs := make(map[string]string) // bulkId -> created resource ID
	errorCount := 0

	var results []map[string]any

	for _, op := range req.Operations {
		if req.FailOnErrors > 0 && errorCount >= req.FailOnErrors {
			break
		}

		data := resolveBulkIDs(op.Data, bulkIDs)
		result := w.executeBulkOp(op, data)

		if op.BulkID != "" {
			result["bulkId"] = op.BulkID
			if loc, _ := result["location"].(string); loc != "" {
				parts := strings.Split(loc, "/")
				if len(parts) > 0 {
					bulkIDs[op.BulkID] = parts[len(parts)-1]
				}
			}
		}

		if status, _ := result["status"].(string); status != "" {
			if parseStatusCode(status) >= 400 {
				errorCount++
			}
		}

		results = append(results, result)
	}

	resp := map[string]any{
		"schemas":    []string{bulkResponseSchema},
		"Operations": results,
	}

	out, _ := json.Marshal(resp)
	rw.WriteHeader(200)
	_, _ = rw.Write(out)
}
