package compliance

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/scim2/test-suite/spec"
)

func TestTR_BinaryRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	raw := []byte("hello, compliance test")
	encoded := base64.StdEncoding.EncodeToString(raw)

	body, resp := createTestResource(t, map[string]any{
		"binaryAttr": encoded,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-BIN / RFC7643-2.3.6-L537: binary must be base64 encoded.
	v, ok := body["binaryAttr"].(string)
	check(t, adv, []string{"EXTRA-TYPE-BIN", "RFC7643-2.3.6-L537"},
		ok && v != "",
		fmt.Sprintf("binaryAttr = %v, want base64 string", body["binaryAttr"]))

	if ok && v != "" {
		_, err := base64.StdEncoding.DecodeString(v)
		check(t, adv, []string{"EXTRA-TYPE-BIN", "RFC7643-2.3.6-L537"},
			err == nil,
			fmt.Sprintf("binaryAttr is not valid base64: %v", err))
	}
}

func TestTR_BooleanRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"booleanAttr": true,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-BOOL: boolean must round-trip as true/false.
	v, ok := body["booleanAttr"].(bool)
	check(t, adv, []string{"EXTRA-TYPE-BOOL"},
		ok && v,
		fmt.Sprintf("booleanAttr = %v (%T), want true", body["booleanAttr"], body["booleanAttr"]))
}

func TestTR_ComplexRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"complexAttr": map[string]any{
			"sub1": "hello",
			"sub2": 99,
			"sub3": true,
		},
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-COMPLEX: complex attribute sub-attributes must round-trip.
	c, ok := body["complexAttr"].(map[string]any)
	check(t, adv, []string{"EXTRA-TYPE-COMPLEX"},
		ok && c != nil,
		fmt.Sprintf("complexAttr = %v, want object", body["complexAttr"]))

	if ok && c != nil {
		s1, _ := c["sub1"].(string)
		check(t, adv, []string{"EXTRA-TYPE-COMPLEX"},
			s1 == "hello",
			fmt.Sprintf("complexAttr.sub1 = %q, want \"hello\"", s1))

		s2, _ := c["sub2"].(float64)
		check(t, adv, []string{"EXTRA-TYPE-COMPLEX"},
			s2 == 99,
			fmt.Sprintf("complexAttr.sub2 = %v, want 99", c["sub2"]))

		// RFC7643-2.3.8-L595: complex must not contain sub-attributes that are complex.
		for k, v := range c {
			_, isMap := v.(map[string]any)
			check(t, adv, []string{"RFC7643-2.3.8-L595"},
				!isMap,
				fmt.Sprintf("complexAttr.%s is a nested complex (not allowed)", k))
		}
	}
}

func TestTR_DateTimeRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"dateTimeAttr": "2025-06-15T10:30:00Z",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-DT: dateTime must round-trip as xsd:dateTime.
	v, ok := body["dateTimeAttr"].(string)
	check(t, adv, []string{"EXTRA-TYPE-DT"},
		ok && strings.Contains(v, "T") && (strings.HasSuffix(v, "Z") || len(v) > 19),
		fmt.Sprintf("dateTimeAttr = %q, want valid xsd:dateTime", v))
}

// -- Extensions --

func TestTR_ExtensionContainer(t *testing.T) {
	adv := requires(t, spec.TestResource)

	// Create with extension attributes under the extension URI namespace.
	body, resp := createTestResource(t, map[string]any{
		"schemas": []string{testResourceSchema, testExtSchema},
		testExtSchema: map[string]any{
			"extString":   "ext-value",
			"extRequired": "required-value",
		},
	})
	// The schemas key in extra overrides the one in createTestResource,
	// so we need to make sure it includes both.
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-EXT-SCHEMAS: schemas array must include extension URI.
	schemas := getStringSlice(body, "schemas")
	check(t, adv, []string{"EXTRA-EXT-SCHEMAS"},
		hasSchema(body, testExtSchema),
		fmt.Sprintf("schemas = %v, missing extension URI %s", schemas, testExtSchema))

	// EXTRA-EXT-CONTAINER / RFC7643-3.3-L981: extension attributes must be
	// namespaced under the extension schema URI.
	extData, ok := body[testExtSchema].(map[string]any)
	check(t, adv, []string{"EXTRA-EXT-CONTAINER", "RFC7643-3.3-L981"},
		ok && extData != nil,
		fmt.Sprintf("extension data not found under key %q", testExtSchema))

	if ok && extData != nil {
		// EXTRA-EXT-ROUNDTRIP: values must round-trip.
		es, _ := extData["extString"].(string)
		check(t, adv, []string{"EXTRA-EXT-ROUNDTRIP"},
			es == "ext-value",
			fmt.Sprintf("extString = %q, want \"ext-value\"", es))

		er, _ := extData["extRequired"].(string)
		check(t, adv, []string{"EXTRA-EXT-ROUNDTRIP"},
			er == "required-value",
			fmt.Sprintf("extRequired = %q, want \"required-value\"", er))
	}

	// Verify on GET too.
	id := idOf(body)
	getResp, err := client.Get("/TestResources/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		extGet, ok := getResp.Body[testExtSchema].(map[string]any)
		check(t, adv, []string{"EXTRA-EXT-ROUNDTRIP"},
			ok && extGet != nil,
			"GET: extension data missing after creation")
	}
}

// -- Case sensitivity in filters --

func TestTR_FilterCaseExact(t *testing.T) {
	adv := requires(t, spec.TestResource, spec.Filter)

	body, resp := createTestResource(t, map[string]any{
		"caseExactString": "CaSeExAcT",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	_ = body

	// EXTRA-CASE-EXACT: filter on caseExact=true should be case-sensitive.
	// Search with wrong case should NOT match.
	filter := "caseExactString eq \"caseexact\""
	qResp, err := client.Get("/TestResources?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	if qResp.StatusCode == 200 {
		tr, _ := qResp.Body["totalResults"].(float64)
		check(t, adv, []string{"EXTRA-CASE-EXACT", "RFC7643-7-L1748"},
			tr == 0,
			fmt.Sprintf("caseExact filter with wrong case matched %v results, want 0", tr))
	}

	// Search with correct case should match.
	filter = "caseExactString eq \"CaSeExAcT\""
	qResp2, err := client.Get("/TestResources?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	if qResp2.StatusCode == 200 {
		tr, _ := qResp2.Body["totalResults"].(float64)
		check(t, adv, []string{"EXTRA-CASE-EXACT", "RFC7643-7-L1748"},
			tr >= 1,
			fmt.Sprintf("caseExact filter with correct case matched %v results, want >= 1", tr))
	}
}

func TestTR_FilterCaseInsensitive(t *testing.T) {
	adv := requires(t, spec.TestResource, spec.Filter)

	body, resp := createTestResource(t, map[string]any{
		"caseInsensitiveString": "MiXeDcAsE",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	_ = body

	// EXTRA-CASE-INSENSITIVE: filter on caseExact=false should match regardless of case.
	filter := "caseInsensitiveString eq \"mixedcase\""
	qResp, err := client.Get("/TestResources?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	if qResp.StatusCode == 200 {
		tr, _ := qResp.Body["totalResults"].(float64)
		check(t, adv, []string{"EXTRA-CASE-INSENSITIVE"},
			tr >= 1,
			fmt.Sprintf("case-insensitive filter matched %v results, want >= 1", tr))
	}
}

// -- Mutability --

func TestTR_ImmutableAttr(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"immutableAttr": "set-at-creation",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)
	identifier, _ := body["identifier"].(string)

	// Value should be set at creation.
	v, _ := body["immutableAttr"].(string)
	check(t, adv, []string{"EXTRA-MUT-IMMUTABLE"},
		v == "set-at-creation",
		fmt.Sprintf("immutableAttr = %q after create, want \"set-at-creation\"", v))

	// PUT with a different value should fail or be ignored.
	putResp, err := client.Put("/TestResources/"+id, map[string]any{
		"schemas":       []string{testResourceSchema},
		"identifier":    identifier,
		"immutableAttr": "changed-value",
	})
	requireOK(t, err)

	if putResp.StatusCode == 200 {
		updated, _ := putResp.Body["immutableAttr"].(string)
		check(t, adv, []string{"EXTRA-MUT-IMMUTABLE"},
			updated == "set-at-creation",
			fmt.Sprintf("immutableAttr changed to %q on PUT, should stay immutable", updated))
	} else {
		check(t, adv, []string{"EXTRA-MUT-IMMUTABLE"},
			putResp.StatusCode == 400,
			fmtStatus("PUT with changed immutable attr", putResp.StatusCode, 400))
	}
}

// -- Attribute type round-trips --

func TestTR_IntegerRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"integerAttr": 42,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /TestResources returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-INT: integer must round-trip without fractional parts.
	v, ok := body["integerAttr"].(float64)
	check(t, adv, []string{"EXTRA-TYPE-INT"},
		ok && v == 42,
		fmt.Sprintf("integerAttr = %v, want 42", body["integerAttr"]))

	// RFC7643-2.3.4-L521: integer must not contain fractional or exponent parts.
	// Verify via GET that the value is a whole number.
	id := idOf(body)
	getResp, err := client.Get("/TestResources/" + id)
	requireOK(t, err)
	if getResp.StatusCode == 200 {
		gv, _ := getResp.Body["integerAttr"].(float64)
		check(t, adv, []string{"RFC7643-2.3.4-L521"},
			gv == float64(int(gv)),
			fmt.Sprintf("integerAttr on GET = %v, contains fractional part", gv))
	}
}

func TestTR_MultiComplexPrimary(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"multiComplex": []map[string]any{
			{"value": "one", "type": "a", "primary": true},
			{"value": "two", "type": "b", "primary": true},
		},
	})
	if resp.StatusCode == 201 {
		// EXTRA-MV-COMPLEX-PRIMARY: at most one primary=true.
		mc, _ := body["multiComplex"].([]any)
		primaryCount := 0
		for _, item := range mc {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if p, _ := m["primary"].(bool); p {
				primaryCount++
			}
		}
		check(t, adv, []string{"EXTRA-MV-COMPLEX-PRIMARY"},
			primaryCount <= 1,
			fmt.Sprintf("multiComplex has %d primary=true, want <= 1", primaryCount))
	} else if resp.StatusCode == 400 {
		// Rejecting duplicate primary is also acceptable.
		check(t, adv, []string{"EXTRA-MV-COMPLEX-PRIMARY"}, true, "")
	} else {
		check(t, adv, []string{"EXTRA-MV-COMPLEX-PRIMARY"}, false,
			fmtStatus("POST /TestResources with duplicate primary", resp.StatusCode, 400))
	}
}

// -- Multi-valued --

func TestTR_MultiValuedRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"multiStrings": []string{"alpha", "beta", "gamma"},
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-MV-ROUNDTRIP: all values must be preserved.
	arr, ok := body["multiStrings"].([]any)
	check(t, adv, []string{"EXTRA-MV-ROUNDTRIP"},
		ok && len(arr) == 3,
		fmt.Sprintf("multiStrings = %v, want 3 elements", body["multiStrings"]))
}

func TestTR_ReadOnlyIgnored(t *testing.T) {
	adv := requires(t, spec.TestResource)

	// Send a readOnly attribute in the request; server should ignore it.
	body, resp := createTestResource(t, map[string]any{
		"readOnlyAttr": "client-value",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-MUT-READONLY: client-provided value must be ignored.
	v, _ := body["readOnlyAttr"].(string)
	check(t, adv, []string{"EXTRA-MUT-READONLY"},
		v != "client-value",
		fmt.Sprintf("readOnlyAttr = %q, server should have ignored client value", v))
}

func TestTR_ReferenceRoundTrip(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"referenceAttr": "https://example.com/resource/123",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-TYPE-REF / RFC7643-2.3.7-L552: reference must be a URI.
	v, ok := body["referenceAttr"].(string)
	check(t, adv, []string{"EXTRA-TYPE-REF", "RFC7643-2.3.7-L552"},
		ok && v != "",
		fmt.Sprintf("referenceAttr = %v, want URI string", body["referenceAttr"]))

	if ok && v != "" {
		_, err := url.Parse(v)
		check(t, adv, []string{"EXTRA-TYPE-REF", "RFC7643-2.3.7-L552"},
			err == nil,
			fmt.Sprintf("referenceAttr is not a valid URI: %v", err))
	}
}

// -- Required extension --

func TestTR_RequiredExtension(t *testing.T) {
	adv := requires(t, spec.TestResource)

	// RFC7643-6-L1655: if extension is required, resource must include it.
	// Try to create a TestResource WITHOUT the required extension.
	resp, err := client.Post("/TestResources", map[string]any{
		"schemas":    []string{testResourceSchema},
		"identifier": "scim-test-noext-" + randomSuffix(),
	})
	requireOK(t, err)

	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/TestResources/" + id) })
		}
	}

	check(t, adv, []string{"RFC7643-6-L1655"},
		resp.StatusCode == 400,
		fmtStatus("POST /TestResources without required extension", resp.StatusCode, 400))
}

// -- Returned behavior --

func TestTR_ReturnedAlways(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"alwaysReturnedAttr": "always-here",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// EXTRA-RET-ALWAYS: must appear even when excludedAttributes is used.
	getResp, err := client.Get("/TestResources/" + id + "?excludedAttributes=alwaysReturnedAttr")
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		v, has := getResp.Body["alwaysReturnedAttr"]
		check(t, adv, []string{"EXTRA-RET-ALWAYS"},
			has && v != nil,
			"returned:always attribute was excluded by excludedAttributes")
	}
}

func TestTR_WriteOnlyAttr(t *testing.T) {
	adv := requires(t, spec.TestResource)

	body, resp := createTestResource(t, map[string]any{
		"writeOnlyAttr": "secret-value",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// EXTRA-MUT-WRITEONLY: writeOnly must be accepted but never returned.
	_, hasWO := body["writeOnlyAttr"]
	check(t, adv, []string{"EXTRA-MUT-WRITEONLY", "EXTRA-RET-NEVER"},
		!hasWO,
		"writeOnlyAttr was returned in POST response")

	// Also verify on GET.
	id := idOf(body)
	getResp, err := client.Get("/TestResources/" + id)
	requireOK(t, err)
	if getResp.StatusCode == 200 {
		_, hasWO = getResp.Body["writeOnlyAttr"]
		check(t, adv, []string{"EXTRA-MUT-WRITEONLY", "EXTRA-RET-NEVER"},
			!hasWO,
			"writeOnlyAttr was returned in GET response")
	}
}
