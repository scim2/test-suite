package spec

import (
	"fmt"
	"net/url"
	"strings"
)

var l7644_3_4 = rfcLines(rfc7644Sections, "12-section-3.4.txt")

var rfc7644_3_4 = []Requirement{
	// Section 3.4.2

	{
		ID:      "RFC7644-3.4.2-L817",
		Level:   Must,
		Summary: "Query responses must use ListResponse schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: l7644_3_4(101),
			EndLine:   l7644_3_4(102),
			EndCol:    55,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "list_response_schema",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Users")
				r.RequireOK(err)

				r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
					"GET /Users: missing ListResponse schema")

				_, hasTR := resp.Body["totalResults"]
				r.Check(hasTR,
					"GET /Users: missing totalResults in ListResponse")
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2-L847",
		Level:   Shall,
		Summary: "Empty query result returns 200 with totalResults 0",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: l7644_3_4(131),
			EndLine:   l7644_3_4(132),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "empty_filter_returns_200",
			Fn: func(r *Run) {
				filter := fmt.Sprintf("userName eq %q", "scim-nonexistent-"+RandomSuffix())
				resp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /Users?filter=... (no match)", resp.StatusCode, 200))

				if resp.Body != nil {
					tr, _ := resp.Body["totalResults"].(float64)
					r.Check(tr == 0,
						fmt.Sprintf("GET /Users?filter=... (no match): totalResults = %v, want 0", resp.Body["totalResults"]))
				}
			},
		}},
	},

	// Section 3.4.2.1

	{
		ID:      "RFC7644-3.4.2.1-L912",
		Level:   Shall,
		Summary: "Too many results shall return 400 with scimType tooMany",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: l7644_3_4(194),
			EndLine:   l7644_3_4(198),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "too_many_results",
			Fn: func(r *Run) {
				spcResp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)
				filterCfg, _ := spcResp.Body["filter"].(map[string]any)
				maxResults, _ := filterCfg["maxResults"].(float64)
				if maxResults == 0 || maxResults > 100 {
					r.Skipf("maxResults=%v, skipping (too large to create that many resources)", maxResults)
				}

				limit := int(maxResults)
				prefix := "scim-test-toomany-" + RandomSuffix()
				for i := range limit + 1 {
					resp, err := r.Client.Post("/Users", map[string]any{
						"schemas":  []string{UserSchema},
						"userName": fmt.Sprintf("%s-%d", prefix, i),
					})
					r.RequireOK(err)
					if resp.StatusCode == 201 {
						if id := IDOf(resp.Body); id != "" {
							r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
						}
					}
				}

				filter := fmt.Sprintf("userName sw %q", prefix)
				qResp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				if qResp.StatusCode == 400 {
					if qResp.Body != nil {
						scimType, _ := qResp.Body["scimType"].(string)
						r.Check(scimType == "tooMany",
							fmt.Sprintf("scimType = %q, want \"tooMany\"", scimType))
					}
				} else {
					resources, _ := qResp.Body["Resources"].([]any)
					r.Check(len(resources) <= limit,
						fmt.Sprintf("returned %d resources, maxResults=%d", len(resources), limit))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.1-L918",
		Level:   Must,
		Summary: "Multi-resource-type queries must process filters per Section 3.4.2.2",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: l7644_3_4(200),
			EndLine:   l7644_3_4(202),
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "multi_resource_query",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/?filter=" + url.QueryEscape("id pr"))
				r.RequireOK(err)

				switch resp.StatusCode {
				case 200:
					r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
						"root query: missing ListResponse schema")
				case 404, 501:
					// Server doesn't support root-level queries.
					r.Logf("root-level query not supported (status %d)", resp.StatusCode)
				default:
					r.Check(false,
						FmtStatus("GET /?filter=id pr", resp.StatusCode, 200))
				}
			},
		}},
	},

	// Section 3.4.2.2

	{
		ID:      "RFC7644-3.4.2.2-L944",
		Level:   Must,
		Summary: "Filter parameter must contain at least one valid expression",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(228),
			EndLine:   l7644_3_4(231),
			EndCol:    53,
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "eq_match",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}
				userName, _ := body["userName"].(string)

				filter := fmt.Sprintf("userName eq %q", userName)
				qResp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(qResp.StatusCode == 200,
					FmtStatus("GET /Users?filter=userName eq ...", qResp.StatusCode, 200))

				tr, _ := qResp.Body["totalResults"].(float64)
				r.Check(tr >= 1,
					fmt.Sprintf("filter eq match: totalResults=%v, want >= 1", tr))
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.2-L1001",
		Level:   Shall,
		Summary: "Boolean/Binary gt comparison shall fail with 400 invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(285),
			EndLine:   l7644_3_4(288),
			EndCol:    46,
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "boolean_gt_invalid",
			Fn: func(r *Run) {
				filter := "active gt true"
				resp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(resp.StatusCode == 400,
					FmtStatus("GET /Users?filter=active gt true", resp.StatusCode, 400))

				if resp.StatusCode == 400 && resp.Body != nil {
					scimType, _ := resp.Body["scimType"].(string)
					r.Check(scimType == "invalidFilter",
						fmt.Sprintf("scimType = %q, want \"invalidFilter\"", scimType))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.2-L1127",
		Level:   Must,
		Summary: "SCIM filters must conform to the ABNF rules",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(411),
			EndLine:   l7644_3_4(412),
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "invalid_syntax",
			Fn: func(r *Run) {
				filter := "((("
				resp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(resp.StatusCode == 400,
					FmtStatus("GET /Users?filter=(((", resp.StatusCode, 400))
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.2-L1197",
		Level:   Must,
		Summary: "Complex attributes must use fully qualified sub-attribute in filters",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(481),
			EndLine:   l7644_3_4(483),
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "sub_attribute",
			Fn: func(r *Run) {
				_, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				filter := "name.givenName pr"
				qResp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(qResp.StatusCode == 200 || qResp.StatusCode == 400,
					FmtStatus("GET /Users?filter=name.givenName pr", qResp.StatusCode, 200))
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.2-L1208",
		Level:   Must,
		Summary: "Unrecognized filter operations must return 400 with invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(492),
			EndLine:   l7644_3_4(498),
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "invalid_operator",
			Fn: func(r *Run) {
				filter := "userName regex \"test\""
				resp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				r.Check(resp.StatusCode == 400,
					FmtStatus("GET /Users?filter=userName regex ...", resp.StatusCode, 400))

				if resp.Body != nil {
					scimType, _ := resp.Body["scimType"].(string)
					r.Check(scimType == "invalidFilter",
						fmt.Sprintf("scimType = %q, want \"invalidFilter\"", scimType))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.2-L1217",
		Level:   Shall,
		Summary: "String filter comparisons shall respect caseExact characteristic",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(500),
			EndLine:   l7644_3_4(502),
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "case_insensitive",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}
				userName, _ := body["userName"].(string)

				upper := strings.ToUpper(userName)
				filter := fmt.Sprintf("userName eq %q", upper)
				qResp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				if qResp.StatusCode == 200 {
					tr, _ := qResp.Body["totalResults"].(float64)
					r.Check(tr >= 1,
						fmt.Sprintf("case-insensitive filter: totalResults=%v, want >= 1 (searched %q for %q)", tr, upper, userName))
				}
			},
		}},
	},

	// Section 3.4.2.3

	{
		ID:      "RFC7644-3.4.2.3-L1319",
		Level:   Shall,
		Summary: "sortOrder shall default to ascending when sortBy is provided",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: l7644_3_4(600),
			EndLine:   l7644_3_4(603),
			EndCol:    33,
		},
		Feature:  Sort,
		Testable: true,
		Tests: []Test{{
			Name: "default_ascending",
			Fn: func(r *Run) {
				for _, name := range []string{"alpha", "beta", "gamma"} {
					resp, err := r.Client.Post("/Users", map[string]any{
						"schemas":  []string{UserSchema},
						"userName": "scim-test-sort-" + name + "-" + RandomSuffix(),
					})
					r.RequireOK(err)
					if resp.StatusCode == 201 {
						if id := IDOf(resp.Body); id != "" {
							r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
						}
					}
				}

				resp, err := r.Client.Get("/Users?sortBy=userName")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /Users?sortBy=userName", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					resources, _ := resp.Body["Resources"].([]any)
					if len(resources) >= 2 {
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
						r.Check(sorted,
							"sortBy=userName: results not in ascending order")
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.3-L1321",
		Level:   Must,
		Summary: "sortOrder must sort according to attribute type (case-insensitive or case-exact)",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: l7644_3_4(603),
			EndLine:   l7644_3_4(610),
			StartCol:  36,
		},
		Feature:  Sort,
		Testable: true,
		Tests: []Test{{
			Name: "by_attribute_type",
			Fn: func(r *Run) {
				for _, name := range []string{"alice", "Bob", "charlie"} {
					resp, err := r.Client.Post("/Users", map[string]any{
						"schemas":  []string{UserSchema},
						"userName": "scim-test-sorttype-" + name + "-" + RandomSuffix(),
					})
					r.RequireOK(err)
					if resp.StatusCode == 201 {
						if id := IDOf(resp.Body); id != "" {
							r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
						}
					}
				}

				resp, err := r.Client.Get("/Users?sortBy=userName&sortOrder=ascending")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /Users?sortBy=userName&sortOrder=ascending", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					resources, _ := resp.Body["Resources"].([]any)
					r.Check(len(resources) > 0,
						"sortBy with sortOrder returned no resources")
				}
			},
		}},
	},

	// Section 3.4.2.4

	{
		ID:      "RFC7644-3.4.2.4-L1358",
		Level:   Shall,
		Summary: "startIndex value less than 1 shall be interpreted as 1",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: l7644_3_4(641),
			EndLine:   l7644_3_4(643),
			StartCol:  40,
			EndCol:    35,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "start_index_zero",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Users?startIndex=0")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /Users?startIndex=0", resp.StatusCode, 200))
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.4-L1363",
		Level:   MustNot,
		Summary: "Service provider must not return more results than count specifies",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: l7644_3_4(647),
			EndLine:   l7644_3_4(648),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "count_limit",
			Fn: func(r *Run) {
				for range 3 {
					r.CreateUser()
				}

				resp, err := r.Client.Get("/Users?count=1")
				r.RequireOK(err)

				if resp.StatusCode != 200 {
					r.Fatalf("GET /Users?count=1 returned %d", resp.StatusCode)
				}

				resources, _ := resp.Body["Resources"].([]any)
				r.Check(len(resources) <= 1,
					fmt.Sprintf("GET /Users?count=1: got %d resources, want <= 1", len(resources)))
			},
		}},
	},

	// Section 3.4.2.5

	{
		ID:      "RFC7644-3.4.2.5-L1438",
		Level:   Must,
		Summary: "attributes and excludedAttributes parameters must be supported",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: l7644_3_4(721),
			EndLine:   l7644_3_4(723),
			StartCol:  31,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "select_user_name",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(body)

				qResp, err := r.Client.Get("/Users/" + id + "?attributes=userName")
				r.RequireOK(err)

				if qResp.StatusCode != 200 {
					r.Fatalf("GET /Users/%s?attributes=userName returned %d", id, qResp.StatusCode)
				}

				userName, _ := qResp.Body["userName"].(string)
				r.Check(userName != "",
					"GET with attributes=userName: userName missing from response")
			},
		}},
	},
	{
		ID:      "RFC7644-3.4.2.5-L1449",
		Level:   Shall,
		Summary: "excludedAttributes shall have no effect on returned:always attributes",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: l7644_3_4(731),
			EndLine:   l7644_3_4(736),
			EndCol:    54,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "excluded_attributes",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(body)

				qResp, err := r.Client.Get("/Users/" + id + "?excludedAttributes=displayName")
				r.RequireOK(err)

				if qResp.StatusCode != 200 {
					r.Fatalf("GET with excludedAttributes returned %d", qResp.StatusCode)
				}

				r.Check(IDOf(qResp.Body) != "",
					"GET with excludedAttributes: id missing (returned:always should not be affected)")

				schemas := GetStringSlice(qResp.Body, "schemas")
				r.Check(len(schemas) > 0,
					"GET with excludedAttributes: schemas missing (returned:always should not be affected)")
			},
		}},
	},

	// Section 3.4.3

	{
		ID:      "RFC7644-3.4.3-L1476",
		Level:   Must,
		Summary: "POST search requests must use SearchRequest schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.3",
			StartLine: l7644_3_4(760),
			EndLine:   l7644_3_4(762),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "post_search",
			Fn: func(r *Run) {
				r.CreateUser()

				resp, err := r.Client.Post("/Users/.search", map[string]any{
					"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:SearchRequest"},
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("POST /Users/.search", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
						"POST /.search response missing ListResponse schema")
				}
			},
		}},
	},
}
