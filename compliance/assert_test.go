package compliance

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/scim2/test-suite/scim"
	"github.com/scim2/test-suite/spec"
)

const (
	userSchema         = "urn:ietf:params:scim:schemas:core:2.0:User"
	groupSchema        = "urn:ietf:params:scim:schemas:core:2.0:Group"
	testResourceSchema = "urn:ietf:params:scim:schemas:test:2.0:TestResource"
	testExtSchema      = "urn:ietf:params:scim:schemas:extension:test:2.0:TestResource"
)

// check asserts a condition and records the result. When advertised is
// false (feature not advertised), failures are recorded as Warn and
// use t.Log instead of t.Error (no hard fail).
func check(t *testing.T, advertised bool, ids []string, ok bool, msg string) {
	t.Helper()
	if ok {
		record(Result{RequirementIDs: ids, Outcome: Pass})
		return
	}

	if advertised {
		record(Result{RequirementIDs: ids, Outcome: Fail, Message: msg})
		for _, id := range ids {
			if req := spec.ByID(id); req != nil {
				t.Errorf("[%s] %s (RFC %d, Section %s, Lines %d-%d): %s",
					req.Level, id, req.Source.RFC, req.Source.Section,
					req.Source.StartLine, req.Source.EndLine, req.Summary)
			}
		}
		t.Error(msg)
	} else {
		record(Result{RequirementIDs: ids, Outcome: Warn, Message: msg})
		for _, id := range ids {
			if req := spec.ByID(id); req != nil {
				t.Logf("[%s] %s (RFC %d, Section %s, Lines %d-%d): %s",
					req.Level, id, req.Source.RFC, req.Source.Section,
					req.Source.StartLine, req.Source.EndLine, req.Summary)
			}
		}
		t.Logf("(soft) %s", msg)
	}
}

// createGroup creates a Group and registers cleanup to delete it.
func createGroup(t *testing.T, members []map[string]any) (map[string]any, *scim.Response) {
	t.Helper()
	body := map[string]any{
		"schemas":     []string{groupSchema},
		"displayName": "scim-test-group-" + randomSuffix(),
	}
	if len(members) > 0 {
		body["members"] = members
	}
	resp, err := client.Post("/Groups", body)
	requireOK(t, err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		t.Cleanup(func() {
			_, _ = client.Delete("/Groups/" + id)
		})
	}

	return resp.Body, resp
}

// createTestResource creates a TestResource and registers cleanup to delete it.
// The required extension is included by default.
func createTestResource(t *testing.T, extra map[string]any) (map[string]any, *scim.Response) {
	t.Helper()
	body := map[string]any{
		"schemas":    []string{testResourceSchema, testExtSchema},
		"identifier": "scim-test-" + randomSuffix(),
		testExtSchema: map[string]any{
			"extRequired": "default",
		},
	}
	maps.Copy(body, extra)
	resp, err := client.Post("/TestResources", body)
	requireOK(t, err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		t.Cleanup(func() {
			_, _ = client.Delete("/TestResources/" + id)
		})
	}

	return resp.Body, resp
}

// createUser creates a test user and registers cleanup to delete it.
func createUser(t *testing.T) (map[string]any, *scim.Response) {
	t.Helper()
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-" + randomSuffix(),
	})
	requireOK(t, err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		t.Cleanup(func() {
			_, _ = client.Delete("/Users/" + id)
		})
	}

	return resp.Body, resp
}

// fmtStatus formats a status code mismatch message.
func fmtStatus(what string, got, want int) string {
	return fmt.Sprintf("%s: got status %d, want %d", what, got, want)
}

// generateSelfSignedCertB64 creates a real self-signed X.509 certificate
// and returns it as a base64-encoded DER string.
func generateSelfSignedCertB64(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	return base64.StdEncoding.EncodeToString(der)
}

// getString extracts a nested string from a JSON map.
func getString(m map[string]any, keys ...string) string {
	for i, key := range keys {
		if i == len(keys)-1 {
			v, _ := m[key].(string)
			return v
		}
		next, ok := m[key].(map[string]any)
		if !ok {
			return ""
		}
		m = next
	}
	return ""
}

// getStringSlice extracts a []string from a JSON interface slice.
func getStringSlice(m map[string]any, key string) []string {
	arr, ok := m[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// hasSchema checks if the schemas array contains the given URI.
func hasSchema(body map[string]any, uri string) bool {
	return slices.Contains(getStringSlice(body, "schemas"), uri)
}

// idOf extracts the "id" field from a SCIM resource body.
func idOf(body map[string]any) string {
	id, _ := body["id"].(string)
	return id
}

// isBase64 checks if a string is valid base64.
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// isValidAttrName checks if a string conforms to the SCIM ATTRNAME ABNF:
// ALPHA *(nameChar) where nameChar = "-" / "_" / DIGIT / ALPHA.
func isValidAttrName(s string) bool {
	if len(s) == 0 {
		return false
	}
	c := s[0]
	if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
		return false
	}
	for i := 1; i < len(s); i++ {
		c = s[i]
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') &&
			(c < '0' || c > '9') && c != '-' && c != '_' {
			return false
		}
	}
	return true
}

// randomSuffix returns a short random hex string for unique test data.
func randomSuffix() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func record(r Result) {
	mu.Lock()
	defer mu.Unlock()
	results = append(results, r)
}

// requireOK fails the test immediately if err is non-nil.
func requireOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requires marks the features this test depends on. Returns false if
// a required feature is not available, in which case the test runs in
// "soft" mode: failures are warnings, not errors. The test should
// still attempt to run.
func requires(t *testing.T, feats ...spec.Feature) bool {
	t.Helper()
	for _, f := range feats {
		if !supported(f) {
			t.Logf("feature %q not advertised, running in soft mode", f)
			return false
		}
	}
	return true
}

// supported returns true if the feature is advertised or force-enabled.
func supported(f spec.Feature) bool {
	if force[f] {
		return true
	}
	switch f {
	case spec.Core:
		return true
	case spec.Filter:
		return features.Filter
	case spec.Patch:
		return features.Patch
	case spec.Bulk:
		return features.Bulk
	case spec.Sort:
		return features.Sort
	case spec.ETag:
		return features.ETag
	case spec.ChangePassword:
		return features.ChangePassword
	case spec.Users:
		return features.Users
	case spec.Groups:
		return features.Groups
	case spec.TestResource:
		return features.TestResource
	case spec.Me:
		return false
	}
	return false
}

// Outcome represents the result of a compliance check.
type Outcome int

const (
	Pass Outcome = iota
	Fail
	Warn // feature not advertised but test ran and failed
	Skip
)

// Result links a test outcome to one or more requirement IDs.
type Result struct {
	RequirementIDs []string
	Outcome        Outcome
	Message        string
}
