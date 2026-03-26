package compliance

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/scim2/test-suite/spec"
)

func TestFilter_BooleanGtInvalid(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// RFC7644-3.4.2.2-L1001: Boolean gt comparison shall fail with 400 invalidFilter.
	filter := "active gt true"
	resp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.2-L1001"},
		resp.StatusCode == 400,
		fmtStatus("GET /Users?filter=active gt true", resp.StatusCode, 400))

	if resp.StatusCode == 400 && resp.Body != nil {
		scimType, _ := resp.Body["scimType"].(string)
		check(t, adv, []string{"RFC7644-3.4.2.2-L1001"},
			scimType == "invalidFilter",
			fmt.Sprintf("scimType = %q, want \"invalidFilter\"", scimType))
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	userName, _ := body["userName"].(string)

	// RFC7644-3.4.2.2-L1217: string filter comparisons shall respect caseExact.
	// userName is caseExact=false by default, so case-insensitive matching is expected.
	upper := strings.ToUpper(userName)
	filter := fmt.Sprintf("userName eq %q", upper)
	qResp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	if qResp.StatusCode == 200 {
		tr, _ := qResp.Body["totalResults"].(float64)
		check(t, adv, []string{"RFC7644-3.4.2.2-L1217"},
			tr >= 1,
			fmt.Sprintf("case-insensitive filter: totalResults=%v, want >= 1 (searched %q for %q)", tr, upper, userName))
	}
}

func TestFilter_EqMatch(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	userName, _ := body["userName"].(string)

	filter := fmt.Sprintf("userName eq %q", userName)
	qResp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	// RFC7644-3.4.2.2-L944: filter must contain valid expression.
	check(t, adv, []string{"RFC7644-3.4.2.2-L944"},
		qResp.StatusCode == 200,
		fmtStatus("GET /Users?filter=userName eq ...", qResp.StatusCode, 200))

	tr, _ := qResp.Body["totalResults"].(float64)
	check(t, adv, []string{"RFC7644-3.4.2.2-L944"},
		tr >= 1,
		fmt.Sprintf("filter eq match: totalResults=%v, want >= 1", tr))
}

func TestFilter_InvalidOperator(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// RFC7644-3.4.2.2-L1208: unrecognized filter operation must return 400 invalidFilter.
	filter := "userName regex \"test\""
	resp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.2-L1208"},
		resp.StatusCode == 400,
		fmtStatus("GET /Users?filter=userName regex ...", resp.StatusCode, 400))

	if resp.Body != nil {
		scimType, _ := resp.Body["scimType"].(string)
		check(t, adv, []string{"RFC7644-3.4.2.2-L1208"},
			scimType == "invalidFilter",
			fmt.Sprintf("scimType = %q, want \"invalidFilter\"", scimType))
	}
}

func TestFilter_InvalidSyntax(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// RFC7644-3.4.2.2-L1127: filters must conform to ABNF.
	// A completely invalid filter string.
	filter := "((("
	resp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.2-L1127"},
		resp.StatusCode == 400,
		fmtStatus("GET /Users?filter=(((", resp.StatusCode, 400))
}

func TestFilter_SubAttribute(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// RFC7644-3.4.2.2-L1197: complex attributes must use fully qualified sub-attribute.
	// name.givenName is the correct notation.
	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	_ = body

	filter := "name.givenName pr"
	qResp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	// Should either succeed (200) or return 400 if the attribute isn't filterable.
	// It must NOT return 500.
	check(t, adv, []string{"RFC7644-3.4.2.2-L1197"},
		qResp.StatusCode == 200 || qResp.StatusCode == 400,
		fmtStatus("GET /Users?filter=name.givenName pr", qResp.StatusCode, 200))
}

func TestSort_ByAttributeType(t *testing.T) {
	adv := requires(t, spec.Users, spec.Sort)

	// RFC7644-3.4.2.3-L1321: sortOrder must sort according to attribute type.
	// userName is case-insensitive, so "Alice" and "alice" should sort together.
	for _, name := range []string{"alice", "Bob", "charlie"} {
		resp, err := client.Post("/Users", map[string]any{
			"schemas":  []string{userSchema},
			"userName": "scim-test-sorttype-" + name + "-" + randomSuffix(),
		})
		requireOK(t, err)
		if resp.StatusCode == 201 {
			if id := idOf(resp.Body); id != "" {
				t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
			}
		}
	}

	resp, err := client.Get("/Users?sortBy=userName&sortOrder=ascending")
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.3-L1321"},
		resp.StatusCode == 200,
		fmtStatus("GET /Users?sortBy=userName&sortOrder=ascending", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		resources, _ := resp.Body["Resources"].([]any)
		// Just verify we got resources back and server didn't error.
		check(t, adv, []string{"RFC7644-3.4.2.3-L1321"},
			len(resources) > 0,
			"sortBy with sortOrder returned no resources")
	}
}

func TestSort_DefaultAscending(t *testing.T) {
	adv := requires(t, spec.Users, spec.Sort)

	// Create users with known ordering.
	for _, name := range []string{"alpha", "beta", "gamma"} {
		resp, err := client.Post("/Users", map[string]any{
			"schemas":  []string{userSchema},
			"userName": "scim-test-sort-" + name + "-" + randomSuffix(),
		})
		requireOK(t, err)
		if resp.StatusCode == 201 {
			if id := idOf(resp.Body); id != "" {
				t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
			}
		}
	}

	// RFC7644-3.4.2.3-L1319: sortOrder shall default to ascending when sortBy is provided.
	resp, err := client.Get("/Users?sortBy=userName")
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.3-L1319"},
		resp.StatusCode == 200,
		fmtStatus("GET /Users?sortBy=userName", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		resources, _ := resp.Body["Resources"].([]any)
		if len(resources) >= 2 {
			// Verify ascending order: each userName <= next.
			sorted := true
			for i := 0; i < len(resources)-1; i++ {
				r1, _ := resources[i].(map[string]any)
				r2, _ := resources[i+1].(map[string]any)
				u1, _ := r1["userName"].(string)
				u2, _ := r2["userName"].(string)
				if strings.ToLower(u1) > strings.ToLower(u2) {
					sorted = false
					break
				}
			}
			check(t, adv, []string{"RFC7644-3.4.2.3-L1319"},
				sorted,
				"sortBy=userName: results not in ascending order")
		}
	}
}
