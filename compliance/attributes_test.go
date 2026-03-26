package compliance

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/scim2/test-suite/spec"
)

func TestAttributeNameFormat(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// RFC7643-2.1-L369: attribute names must conform to ATTRNAME ABNF.
	// ATTRNAME = ALPHA *(nameChar) where nameChar = "-" / "_" / DIGIT / ALPHA
	valid := true
	for key := range body {
		if !isValidAttrName(key) {
			valid = false
			t.Logf("invalid attribute name: %q", key)
		}
	}
	check(t, adv, []string{"RFC7643-2.1-L369"},
		valid,
		"response contains attribute names not conforming to ATTRNAME ABNF")
}

func TestAttributes_ExcludedAttributes(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Exclude displayName.
	qResp, err := client.Get("/Users/" + id + "?excludedAttributes=displayName")
	requireOK(t, err)

	if qResp.StatusCode != 200 {
		t.Fatalf("GET with excludedAttributes returned %d", qResp.StatusCode)
	}

	// RFC7644-3.4.2.5-L1449: excludedAttributes shall have no effect on returned:always.
	// id has returned:always, so it must still be present.
	check(t, adv, []string{"RFC7644-3.4.2.5-L1449"},
		idOf(qResp.Body) != "",
		"GET with excludedAttributes: id missing (returned:always should not be affected)")

	schemas := getStringSlice(qResp.Body, "schemas")
	check(t, adv, []string{"RFC7644-3.4.2.5-L1449"},
		len(schemas) > 0,
		"GET with excludedAttributes: schemas missing (returned:always should not be affected)")
}

func TestAttributes_SelectOnList(t *testing.T) {
	adv := requires(t, spec.Users)

	createUser(t)

	// Request list with attributes parameter.
	resp, err := client.Get("/Users?attributes=" + url.QueryEscape("userName"))
	requireOK(t, err)

	if resp.StatusCode != 200 {
		t.Fatalf("GET /Users?attributes=userName returned %d", resp.StatusCode)
	}

	resources, _ := resp.Body["Resources"].([]any)
	if len(resources) == 0 {
		t.Skip("no resources returned")
	}

	first, ok := resources[0].(map[string]any)
	if !ok {
		t.Fatal("first resource is not a JSON object")
	}

	// RFC7644-3.9-L3596: each resource must contain minimum set + requested attrs.
	check(t, adv, []string{"RFC7644-3.9-L3596"},
		idOf(first) != "",
		"list with attributes=userName: resource missing id")

	userName, _ := first["userName"].(string)
	check(t, adv, []string{"RFC7644-3.9-L3596"},
		userName != "",
		"list with attributes=userName: resource missing userName")

	// Non-requested attributes should not be present (except returned:always).
	_, hasDisplayName := first["displayName"]
	check(t, adv, []string{"RFC7644-3.9-L3596"},
		!hasDisplayName,
		fmt.Sprintf("list with attributes=userName: displayName should not be returned, got %v", first["displayName"]))
}

func TestAttributes_SelectUserName(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Request only userName.
	qResp, err := client.Get("/Users/" + id + "?attributes=userName")
	requireOK(t, err)

	if qResp.StatusCode != 200 {
		t.Fatalf("GET /Users/%s?attributes=userName returned %d", id, qResp.StatusCode)
	}

	// RFC7644-3.4.2.5-L1438: attributes parameter must be supported.
	userName, _ := qResp.Body["userName"].(string)
	check(t, adv, []string{"RFC7644-3.4.2.5-L1438"},
		userName != "",
		"GET with attributes=userName: userName missing from response")

	// RFC7644-3.9-L3596: must include minimum set (id, schemas) + requested attributes.
	check(t, adv, []string{"RFC7644-3.9-L3596"},
		idOf(qResp.Body) != "",
		"GET with attributes=userName: id missing (should always be returned)")

	schemas := getStringSlice(qResp.Body, "schemas")
	check(t, adv, []string{"RFC7644-3.9-L3596"},
		len(schemas) > 0,
		"GET with attributes=userName: schemas missing (should always be returned)")
}

func TestCountry_ISO3166(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create user with an address containing a country code.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-country-" + randomSuffix(),
		"addresses": []map[string]any{
			{
				"type":    "work",
				"country": "US",
			},
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode == 201 && resp.Body != nil {
		addrs, _ := resp.Body["addresses"].([]any)
		for _, a := range addrs {
			addr, ok := a.(map[string]any)
			if !ok {
				continue
			}
			country, _ := addr["country"].(string)
			if country == "" {
				continue
			}
			// RFC7643-4.1.2-L1322: country must be ISO 3166-1 alpha-2.
			check(t, adv, []string{"RFC7643-4.1.2-L1322"},
				len(country) == 2 && country[0] >= 'A' && country[0] <= 'Z' && country[1] >= 'A' && country[1] <= 'Z',
				fmt.Sprintf("country = %q, want 2-letter uppercase ISO 3166-1 code", country))
		}
	}
}

func TestEmail_RFC5321Format(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// Create a user with an email to test format.
	id := idOf(body)
	userName, _ := body["userName"].(string)
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
		"emails": []map[string]any{
			{"value": "valid@example.com", "type": "work"},
		},
	})
	requireOK(t, err)

	if putResp.StatusCode == 200 && putResp.Body != nil {
		emails, _ := putResp.Body["emails"].([]any)
		for _, e := range emails {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			v, _ := em["value"].(string)
			// RFC7643-4.1.2-L1244: email values should be per RFC 5321.
			// Basic check: must contain exactly one @.
			atCount := 0
			for _, c := range v {
				if c == '@' {
					atCount++
				}
			}
			check(t, adv, []string{"RFC7643-4.1.2-L1244"},
				atCount == 1 && len(v) > 3,
				fmt.Sprintf("email %q does not look like RFC 5321 format", v))
		}
	}
}

func TestExternalId_NotSetByServer(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create without externalId.
	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// RFC7643-3.1-L890: externalId must not be specified by the service provider.
	_, hasExtId := body["externalId"]
	check(t, adv, []string{"RFC7643-3.1-L890"},
		!hasExtId,
		"POST /Users without externalId: server added externalId to response")
}

func TestIdStable(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Fetch twice.
	resp1, err := client.Get("/Users/" + id)
	requireOK(t, err)
	resp2, err := client.Get("/Users/" + id)
	requireOK(t, err)

	// RFC7643-3.1-L868: id must be stable and non-reassignable.
	check(t, adv, []string{"RFC7643-3.1-L868"},
		idOf(resp1.Body) == id && idOf(resp2.Body) == id,
		fmt.Sprintf("id not stable: got %q and %q, want %q", idOf(resp1.Body), idOf(resp2.Body), id))
}

func TestMultiValued_Canonicalize(t *testing.T) {
	adv := requires(t, spec.Users)

	// RFC7643-2.4-L655: service providers should canonicalize values.
	// Create user with non-canonical email casing.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-canon-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "TEST@EXAMPLE.COM", "type": "work"},
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode == 201 && resp.Body != nil {
		emails, _ := resp.Body["emails"].([]any)
		if len(emails) > 0 {
			first, _ := emails[0].(map[string]any)
			v, _ := first["value"].(string)
			// Canonicalization might lowercase, or preserve -- both are okay
			// as long as it's a valid value. We check it wasn't mangled.
			check(t, adv, []string{"RFC7643-2.4-L655"},
				v != "",
				"email value was empty after canonicalization")
		}
	}
}

func TestMultiValued_NoDuplicateTypeValue(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create user with duplicate (type, value) email entries.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-dedup-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "dup@example.com", "type": "work"},
			{"value": "dup@example.com", "type": "work"},
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode == 201 && resp.Body != nil {
		emails, _ := resp.Body["emails"].([]any)

		// RFC7643-2.4-L663: should not return same (type, value) more than once.
		type tv struct{ t, v string }
		seen := make(map[tv]bool)
		hasDupes := false
		for _, e := range emails {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			key := tv{
				t: fmt.Sprintf("%v", em["type"]),
				v: fmt.Sprintf("%v", em["value"]),
			}
			if seen[key] {
				hasDupes = true
			}
			seen[key] = true
		}
		check(t, adv, []string{"RFC7643-2.4-L663"},
			!hasDupes,
			fmt.Sprintf("emails contains duplicate (type, value) pairs: %v", emails))
	} else if resp.StatusCode == 400 {
		// Server correctly rejected duplicate input.
		check(t, adv, []string{"RFC7643-2.4-L663"}, true, "")
	} else {
		check(t, adv, []string{"RFC7643-2.4-L663"}, false,
			fmtStatus("POST /Users with duplicate (type,value)", resp.StatusCode, 400))
	}
}

func TestPostSearch(t *testing.T) {
	adv := requires(t, spec.Users)

	createUser(t)

	// RFC7644-3.4.3-L1476: POST search must use SearchRequest schema URI.
	resp, err := client.Post("/Users/.search", map[string]any{
		"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:SearchRequest"},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.3-L1476"},
		resp.StatusCode == 200,
		fmtStatus("POST /Users/.search", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		check(t, adv, []string{"RFC7644-3.4.3-L1476"},
			hasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
			"POST /.search response missing ListResponse schema")
	}
}

func TestPrimaryUniqueness(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create user with two emails, both marked primary.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-primary-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "a@example.com", "type": "work", "primary": true},
			{"value": "b@example.com", "type": "home", "primary": true},
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	// RFC7643-2.4-L634: primary=true must appear no more than once.
	if resp.StatusCode == 201 && resp.Body != nil {
		emails, _ := resp.Body["emails"].([]any)
		primaryCount := 0
		for _, e := range emails {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			if p, _ := em["primary"].(bool); p {
				primaryCount++
			}
		}
		check(t, adv, []string{"RFC7643-2.4-L634"},
			primaryCount <= 1,
			fmt.Sprintf("emails has %d primary=true values, want at most 1", primaryCount))
	} else if resp.StatusCode == 400 {
		// Server correctly rejected the input.
		check(t, adv, []string{"RFC7643-2.4-L634"}, true, "")
	} else {
		check(t, adv, []string{"RFC7643-2.4-L634"}, false,
			fmtStatus("POST /Users with duplicate primary", resp.StatusCode, 400))
	}
}

func TestPut_ReplaceValues(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)
	userName, _ := body["userName"].(string)

	newDisplay := "Replaced-" + randomSuffix()
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":     []string{userSchema},
		"userName":    userName,
		"displayName": newDisplay,
	})
	requireOK(t, err)

	if putResp.StatusCode != 200 {
		t.Fatalf("PUT returned %d", putResp.StatusCode)
	}

	// RFC7644-3.5.1-L1644: readWrite values provided SHALL replace existing.
	dn, _ := putResp.Body["displayName"].(string)
	check(t, adv, []string{"RFC7644-3.5.1-L1644"},
		dn == newDisplay,
		fmt.Sprintf("PUT replace: displayName = %q, want %q", dn, newDisplay))

	// RFC7644-3.5.1-L1665: readOnly values shall be ignored.
	check(t, adv, []string{"RFC7644-3.5.1-L1665"},
		idOf(putResp.Body) == id,
		"PUT: id changed (readOnly should be ignored)")
}

func TestResponse_SchemasPresent(t *testing.T) {
	// RFC7643-2.1-L366: resources must be JSON and specify schema.
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	check(t, adv, []string{"RFC7643-2.1-L366"},
		body != nil,
		"POST /Users: response is not a JSON object")

	schemas := getStringSlice(body, "schemas")
	check(t, adv, []string{"RFC7643-2.1-L366"},
		len(schemas) > 0,
		"POST /Users: response missing schemas attribute")
}

func TestUnassignedEquivalence(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create user without displayName.
	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// RFC7643-2.5-L682: unassigned, null, and empty array are equivalent.
	// displayName was not sent, so it should be absent or null.
	dn, hasDN := body["displayName"]
	check(t, adv, []string{"RFC7643-2.5-L682"},
		!hasDN || dn == nil || dn == "",
		fmt.Sprintf("unassigned displayName returned as %v", dn))
}

func TestWriteOnly_NotReturned(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	// RFC7643-7-L1769: writeOnly attribute values shall not be returned.
	_, hasPassword := body["password"]
	check(t, adv, []string{"RFC7643-7-L1769"},
		!hasPassword,
		"writeOnly attribute 'password' was returned in response")

	// Also check on GET.
	id := idOf(body)
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		_, hasPassword = getResp.Body["password"]
		check(t, adv, []string{"RFC7643-7-L1769"},
			!hasPassword,
			"GET returned writeOnly attribute 'password'")
	}
}

func TestX509_Base64DER(t *testing.T) {
	adv := requires(t, spec.Users)

	// Generate a real self-signed X.509 certificate for testing.
	certB64 := generateSelfSignedCertB64(t)

	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-x509-" + randomSuffix(),
		"x509Certificates": []map[string]any{
			{"value": certB64},
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode == 201 && resp.Body != nil {
		certs, _ := resp.Body["x509Certificates"].([]any)
		for _, c := range certs {
			cert, ok := c.(map[string]any)
			if !ok {
				continue
			}
			v, _ := cert["value"].(string)
			// RFC7643-4.1.2-L1378: x509 values must be base64 encoded DER.
			check(t, adv, []string{"RFC7643-4.1.2-L1378"},
				v != "" && isBase64(v),
				fmt.Sprintf("x509Certificate value is not valid base64: %q", v))
		}
	}
}
