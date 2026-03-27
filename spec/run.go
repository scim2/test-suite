package spec

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
	"net/http"
	"slices"
	"time"

	"github.com/scim2/test-suite/scim"
)

const (
	UserSchema         = "urn:ietf:params:scim:schemas:core:2.0:User"
	GroupSchema        = "urn:ietf:params:scim:schemas:core:2.0:Group"
	TestResourceSchema = "urn:ietf:params:scim:schemas:test:2.0:TestResource"
	TestExtSchema      = "urn:ietf:params:scim:schemas:extension:test:2.0:TestResource"
	PatchOpSchema      = "urn:ietf:params:scim:api:messages:2.0:PatchOp"
	BulkSchema         = "urn:ietf:params:scim:api:messages:2.0:BulkRequest"
)

// FmtStatus formats a status code mismatch message.
func FmtStatus(what string, got, want int) string {
	return fmt.Sprintf("%s: got status %d, want %d", what, got, want)
}

// GetString extracts a nested string from a JSON map.
func GetString(m map[string]any, keys ...string) string {
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

// GetStringSlice extracts a []string from a JSON interface slice.
func GetStringSlice(m map[string]any, key string) []string {
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

// HasSchema checks if the schemas array contains the given URI.
func HasSchema(body map[string]any, uri string) bool {
	return slices.Contains(GetStringSlice(body, "schemas"), uri)
}

// IDOf extracts the "id" field from a SCIM resource body.
func IDOf(body map[string]any) string {
	id, _ := body["id"].(string)
	return id
}

// IsBase64 checks if a string is valid base64.
func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// IsValidAttrName checks if a string conforms to the SCIM ATTRNAME ABNF:
// ALPHA *(nameChar) where nameChar = "-" / "_" / DIGIT / ALPHA.
func IsValidAttrName(s string) bool {
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

// RandomSuffix returns a short random hex string for unique test data.
func RandomSuffix() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Run provides assertions, HTTP helpers, and lifecycle management for
// a single compliance test. It carries the SCIM client so that test
// functions defined in spec/ can make HTTP calls without importing
// any other package.
type Run struct {
	Client   *scim.Client
	Features *scim.Features

	failed   bool
	skipped  bool
	msgs     []string
	logs     []string
	cleanups []func()
}

// Check asserts a condition. When ok is false, Errorf is called.
func (r *Run) Check(ok bool, msg string) {
	if !ok {
		r.Errorf("%s", msg)
	}
}

// Cleanup registers a function to be called after the test completes.
func (r *Run) Cleanup(fn func()) {
	r.cleanups = append(r.cleanups, fn)
}

// CreateGroup creates a Group and registers cleanup to delete it.
func (r *Run) CreateGroup(members []map[string]any) (map[string]any, *scim.Response) {
	body := map[string]any{
		"schemas":     []string{GroupSchema},
		"displayName": "scim-test-group-" + RandomSuffix(),
	}
	if len(members) > 0 {
		body["members"] = members
	}
	resp, err := r.Client.Post("/Groups", body)
	r.RequireOK(err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		r.Cleanup(func() { _, _ = r.Client.Delete("/Groups/" + id) })
	}

	return resp.Body, resp
}

// CreateTestResource creates a TestResource and registers cleanup.
func (r *Run) CreateTestResource(extra map[string]any) (map[string]any, *scim.Response) {
	body := map[string]any{
		"schemas":    []string{TestResourceSchema, TestExtSchema},
		"identifier": "scim-test-" + RandomSuffix(),
		TestExtSchema: map[string]any{
			"extRequired": "default",
		},
	}
	maps.Copy(body, extra)
	resp, err := r.Client.Post("/TestResources", body)
	r.RequireOK(err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		r.Cleanup(func() { _, _ = r.Client.Delete("/TestResources/" + id) })
	}

	return resp.Body, resp
}

// CreateUser creates a test user and registers cleanup to delete it.
func (r *Run) CreateUser() (map[string]any, *scim.Response) {
	resp, err := r.Client.Post("/Users", map[string]any{
		"schemas":  []string{UserSchema},
		"userName": "scim-test-" + RandomSuffix(),
	})
	r.RequireOK(err)

	if id, ok := resp.Body["id"].(string); ok && resp.StatusCode == 201 {
		r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
	}

	return resp.Body, resp
}

// DoWithHeaders executes an HTTP request with additional headers.
func (r *Run) DoWithHeaders(method, path string, body map[string]any, headers map[string]string) *scim.Response {
	resp, err := r.Client.DoWithHeaders(method, path, body, headers)
	r.RequireOK(err)
	return resp
}

// Errorf records a test failure with a formatted message.
func (r *Run) Errorf(format string, args ...any) {
	r.failed = true
	r.msgs = append(r.msgs, fmt.Sprintf(format, args...))
}

// Execute runs a test function, recovering from Fatalf/Skipf panics
// and running cleanup functions afterward.
func (r *Run) Execute(fn func(r *Run)) {
	defer r.RunCleanups()
	defer func() {
		if v := recover(); v != nil {
			if _, ok := v.(fatalSentinel); !ok {
				panic(v)
			}
		}
	}()
	fn(r)
}

// Failed reports whether the test has failed.
func (r *Run) Failed() bool { return r.failed }

// Fatalf records a test failure and stops execution immediately.
func (r *Run) Fatalf(format string, args ...any) {
	r.failed = true
	r.msgs = append(r.msgs, fmt.Sprintf(format, args...))
	panic(fatalSentinel{})
}

// GenerateSelfSignedCertB64 creates a self-signed X.509 certificate
// and returns it as a base64-encoded DER string.
func (r *Run) GenerateSelfSignedCertB64() string {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		r.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		r.Fatalf("create certificate: %v", err)
	}
	return base64.StdEncoding.EncodeToString(der)
}

// Logf records an informational message.
func (r *Run) Logf(format string, args ...any) {
	r.logs = append(r.logs, fmt.Sprintf(format, args...))
}

// Logs returns all informational messages recorded during the test.
func (r *Run) Logs() []string { return r.logs }

// Messages returns all error messages recorded during the test.
func (r *Run) Messages() []string { return r.msgs }

// RawClient returns a Client with no authentication, for testing
// unauthenticated access.
func (r *Run) RawClient() *scim.Client {
	return &scim.Client{
		BaseURL:    r.Client.BaseURL,
		HTTPClient: &http.Client{},
	}
}

// RequireOK fails the test immediately if err is non-nil.
func (r *Run) RequireOK(err error) {
	if err != nil {
		r.Fatalf("unexpected error: %v", err)
	}
}

// RunCleanups executes all registered cleanup functions in LIFO order.
func (r *Run) RunCleanups() {
	for i := len(r.cleanups) - 1; i >= 0; i-- {
		r.cleanups[i]()
	}
}

// Skipf marks the test as skipped and stops execution.
func (r *Run) Skipf(format string, args ...any) {
	r.skipped = true
	r.logs = append(r.logs, fmt.Sprintf(format, args...))
	panic(fatalSentinel{})
}

// Skipped reports whether the test was skipped.
func (r *Run) Skipped() bool { return r.skipped }

// fatalSentinel is used by Fatalf to stop test execution via panic/recover.
type fatalSentinel struct{}
