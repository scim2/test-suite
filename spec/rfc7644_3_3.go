package spec

import "fmt"

var l7644_3_3 = rfcLines(rfc7644Sections, "11-section-3.3.txt")

var rfc7644_3_3 = []Requirement{
	// Section 3.3

	{
		ID:      "RFC7644-3.3-L584",
		Level:   Shall,
		Summary: "readOnly attributes in POST request body shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(10),
			EndLine:   l7644_3_3(11),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "readonly_ignored_on_create",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-ro-" + RandomSuffix(),
					"id":       "client-specified-id",
					"meta": map[string]any{
						"resourceType": "Invalid",
						"created":      "1999-01-01T00:00:00Z",
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode != 201 {
					r.Skipf("POST /Users returned %d, cannot verify readOnly handling", resp.StatusCode)
				}

				r.Check(IDOf(resp.Body) != "client-specified-id",
					"POST /Users: server used client-specified id instead of generating one")
			},
		}},
	},
	{
		ID:      "RFC7644-3.3-L602",
		Level:   Shall,
		Summary: "Successful resource creation returns HTTP 201 Created",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(28),
			EndLine:   l7644_3_3(29),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "create_returns_201",
			Fn: func(r *Run) {
				_, resp := r.CreateUser()

				r.Check(resp.StatusCode == 201,
					FmtStatus("POST /Users", resp.StatusCode, 201))
			},
		}},
	},
	{
		ID:      "RFC7644-3.3-L603",
		Level:   Should,
		Summary: "Response body should contain the newly created resource representation",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(30),
			EndLine:   l7644_3_3(31),
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "create_response_body",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()

				r.Check(resp.StatusCode == 201 && body != nil && IDOf(body) != "",
					"POST /Users: response body should contain the created resource representation")
			},
		}},
	},
	{
		ID:      "RFC7644-3.3-L605",
		Level:   Shall,
		Summary: "Created resource URI shall be in Location header and meta.location",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(31),
			EndLine:   l7644_3_3(37),
			StartCol:  51,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "location_header_and_meta",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				loc := resp.Header.Get("Location")
				r.Check(loc != "",
					"POST /Users: missing Location header")

				metaLoc := GetString(body, "meta", "location")
				r.Check(metaLoc != "",
					"POST /Users: missing meta.location in response body")
			},
		}},
	},
	{
		ID:      "RFC7644-3.3-L625",
		Level:   Must,
		Summary: "Duplicate resource creation must return 409 Conflict with uniqueness scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(50),
			EndLine:   l7644_3_3(54),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "duplicate_returns_409",
			Fn: func(r *Run) {
				userName := "scim-test-dup-" + RandomSuffix()

				resp1, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
				})
				r.RequireOK(err)
				if id := IDOf(resp1.Body); id != "" {
					r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
				}

				resp2, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
				})
				r.RequireOK(err)
				if resp2.StatusCode == 201 {
					if id := IDOf(resp2.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				r.Check(resp2.StatusCode == 409,
					FmtStatus("POST /Users (duplicate)", resp2.StatusCode, 409))

				if resp2.Body != nil {
					scimType, _ := resp2.Body["scimType"].(string)
					r.Check(scimType == "uniqueness",
						fmt.Sprintf("POST /Users (duplicate): scimType = %q, want \"uniqueness\"", scimType))
				}
			},
		}},
	},

	// Section 3.3.1

	{
		ID:      "RFC7644-3.3.1-L712",
		Level:   Shall,
		Summary: "meta.resourceType shall be set by service provider to match endpoint",
		Source: Source{
			RFC:       7644,
			Section:   "3.3.1",
			StartLine: l7644_3_3(138),
			EndLine:   l7644_3_3(142),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "meta_resource_type_matches_endpoint",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				r.Check(GetString(body, "meta", "resourceType") == "User",
					fmt.Sprintf("POST /Users: meta.resourceType = %q, want \"User\"",
						GetString(body, "meta", "resourceType")))
			},
		}},
	},
}
