package spec

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Extra contains additional compliance checks that go beyond the
// literal RFC requirements. These exercise attribute type combinations,
// extension handling, and edge cases that the standard User/Group
// schemas do not cover.
var Extra = []Requirement{
	// Attribute types.
	{
		ID:       "EXTRA-TYPE-INT",
		Level:    Must,
		Summary:  "Integer attributes must round-trip without fractional parts",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "integer_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"integerAttr": 42,
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST /TestResources returned %d", resp.StatusCode)
					}

					v, ok := body["integerAttr"].(float64)
					r.Check(ok && v == 42,
						fmt.Sprintf("integerAttr = %v, want 42", body["integerAttr"]))

					id := IDOf(body)
					getResp, err := r.Client.Get("/TestResources/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						gv, _ := getResp.Body["integerAttr"].(float64)
						r.Check(gv == float64(int(gv)),
							fmt.Sprintf("integerAttr on GET = %v, contains fractional part", gv))
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-TYPE-BOOL",
		Level:    Must,
		Summary:  "Boolean attributes must round-trip as JSON true/false",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "boolean_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"booleanAttr": true,
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					v, ok := body["booleanAttr"].(bool)
					r.Check(ok && v,
						fmt.Sprintf("booleanAttr = %v (%T), want true", body["booleanAttr"], body["booleanAttr"]))
				},
			},
		},
	},
	{
		ID:       "EXTRA-TYPE-DT",
		Level:    Must,
		Summary:  "DateTime attributes must round-trip as valid xsd:dateTime strings",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "datetime_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"dateTimeAttr": "2025-06-15T10:30:00Z",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					v, ok := body["dateTimeAttr"].(string)
					r.Check(ok && strings.Contains(v, "T") && (strings.HasSuffix(v, "Z") || len(v) > 19),
						fmt.Sprintf("dateTimeAttr = %q, want valid xsd:dateTime", v))
				},
			},
		},
	},
	{
		ID:       "EXTRA-TYPE-BIN",
		Level:    Must,
		Summary:  "Binary attributes must round-trip as base64-encoded strings",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "binary_round_trip",
				Fn: func(r *Run) {
					raw := []byte("hello, compliance test")
					encoded := base64.StdEncoding.EncodeToString(raw)

					body, resp := r.CreateTestResource(map[string]any{
						"binaryAttr": encoded,
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					v, ok := body["binaryAttr"].(string)
					r.Check(ok && v != "",
						fmt.Sprintf("binaryAttr = %v, want base64 string", body["binaryAttr"]))

					if ok && v != "" {
						_, err := base64.StdEncoding.DecodeString(v)
						r.Check(err == nil,
							fmt.Sprintf("binaryAttr is not valid base64: %v", err))
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-TYPE-REF",
		Level:    Must,
		Summary:  "Reference attributes must round-trip as URI strings",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "reference_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"referenceAttr": "https://example.com/resource/123",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					v, ok := body["referenceAttr"].(string)
					r.Check(ok && v != "",
						fmt.Sprintf("referenceAttr = %v, want URI string", body["referenceAttr"]))

					if ok && v != "" {
						_, err := url.Parse(v)
						r.Check(err == nil,
							fmt.Sprintf("referenceAttr is not a valid URI: %v", err))
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-TYPE-COMPLEX",
		Level:    Must,
		Summary:  "Complex attributes must round-trip with sub-attributes intact",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "complex_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"complexAttr": map[string]any{
							"sub1": "hello",
							"sub2": 99,
							"sub3": true,
						},
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					c, ok := body["complexAttr"].(map[string]any)
					r.Check(ok && c != nil,
						fmt.Sprintf("complexAttr = %v, want object", body["complexAttr"]))

					if ok && c != nil {
						s1, _ := c["sub1"].(string)
						r.Check(s1 == "hello",
							fmt.Sprintf("complexAttr.sub1 = %q, want \"hello\"", s1))

						s2, _ := c["sub2"].(float64)
						r.Check(s2 == 99,
							fmt.Sprintf("complexAttr.sub2 = %v, want 99", c["sub2"]))
					}
				},
			},
		},
	},

	// Mutability.
	{
		ID:       "EXTRA-MUT-IMMUTABLE",
		Level:    Must,
		Summary:  "Immutable attributes must be accepted at creation and rejected on update",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "immutable_attr",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"immutableAttr": "set-at-creation",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}
					id := IDOf(body)
					identifier, _ := body["identifier"].(string)

					v, _ := body["immutableAttr"].(string)
					r.Check(v == "set-at-creation",
						fmt.Sprintf("immutableAttr = %q after create, want \"set-at-creation\"", v))

					putResp, err := r.Client.Put("/TestResources/"+id, map[string]any{
						"schemas":       []string{TestResourceSchema},
						"identifier":    identifier,
						"immutableAttr": "changed-value",
					})
					r.RequireOK(err)

					if putResp.StatusCode == 200 {
						updated, _ := putResp.Body["immutableAttr"].(string)
						r.Check(updated == "set-at-creation",
							fmt.Sprintf("immutableAttr changed to %q on PUT, should stay immutable", updated))
					} else {
						r.Check(putResp.StatusCode == 400,
							FmtStatus("PUT with changed immutable attr", putResp.StatusCode, 400))
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-MUT-WRITEONLY",
		Level:    Must,
		Summary:  "WriteOnly attributes must be accepted but never returned",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "write_only_attr",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"writeOnlyAttr": "secret-value",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					_, hasWO := body["writeOnlyAttr"]
					r.Check(!hasWO,
						"writeOnlyAttr was returned in POST response")

					id := IDOf(body)
					getResp, err := r.Client.Get("/TestResources/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						_, hasWO = getResp.Body["writeOnlyAttr"]
						r.Check(!hasWO,
							"writeOnlyAttr was returned in GET response")
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-MUT-READONLY",
		Level:    Must,
		Summary:  "ReadOnly attributes provided by the client must be ignored",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "read_only_ignored",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"readOnlyAttr": "client-value",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					v, _ := body["readOnlyAttr"].(string)
					r.Check(v != "client-value",
						fmt.Sprintf("readOnlyAttr = %q, server should have ignored client value", v))
				},
			},
		},
	},

	// Returned behavior.
	{
		ID:       "EXTRA-RET-ALWAYS",
		Level:    Must,
		Summary:  "Returned:always attributes must appear even when excluded",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "returned_always",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"alwaysReturnedAttr": "always-here",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}
					id := IDOf(body)

					getResp, err := r.Client.Get("/TestResources/" + id + "?excludedAttributes=alwaysReturnedAttr")
					r.RequireOK(err)

					if getResp.StatusCode == 200 {
						v, has := getResp.Body["alwaysReturnedAttr"]
						r.Check(has && v != nil,
							"returned:always attribute was excluded by excludedAttributes")
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-RET-NEVER",
		Level:    Must,
		Summary:  "Returned:never attributes must never appear in responses",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "returned_never",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"writeOnlyAttr": "secret-value",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					_, hasWO := body["writeOnlyAttr"]
					r.Check(!hasWO,
						"writeOnlyAttr was returned in POST response")

					id := IDOf(body)
					getResp, err := r.Client.Get("/TestResources/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						_, hasWO = getResp.Body["writeOnlyAttr"]
						r.Check(!hasWO,
							"writeOnlyAttr was returned in GET response")
					}
				},
			},
		},
	},

	// Case sensitivity.
	{
		ID:       "EXTRA-CASE-EXACT",
		Level:    Must,
		Summary:  "Filters on caseExact=true attributes must match case-sensitively",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "filter_case_exact",
				Fn: func(r *Run) {
					_, resp := r.CreateTestResource(map[string]any{
						"caseExactString": "CaSeExAcT",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					// Wrong case should NOT match.
					filter := "caseExactString eq \"caseexact\""
					qResp, err := r.Client.Get("/TestResources?filter=" + url.QueryEscape(filter))
					r.RequireOK(err)

					if qResp.StatusCode == 200 {
						tr, _ := qResp.Body["totalResults"].(float64)
						r.Check(tr == 0,
							fmt.Sprintf("caseExact filter with wrong case matched %v results, want 0", tr))
					}

					// Correct case should match.
					filter = "caseExactString eq \"CaSeExAcT\""
					qResp2, err := r.Client.Get("/TestResources?filter=" + url.QueryEscape(filter))
					r.RequireOK(err)

					if qResp2.StatusCode == 200 {
						tr, _ := qResp2.Body["totalResults"].(float64)
						r.Check(tr >= 1,
							fmt.Sprintf("caseExact filter with correct case matched %v results, want >= 1", tr))
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-CASE-INSENSITIVE",
		Level:    Must,
		Summary:  "Filters on caseExact=false attributes must match case-insensitively",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "filter_case_insensitive",
				Fn: func(r *Run) {
					_, resp := r.CreateTestResource(map[string]any{
						"caseInsensitiveString": "MiXeDcAsE",
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					filter := "caseInsensitiveString eq \"mixedcase\""
					qResp, err := r.Client.Get("/TestResources?filter=" + url.QueryEscape(filter))
					r.RequireOK(err)

					if qResp.StatusCode == 200 {
						tr, _ := qResp.Body["totalResults"].(float64)
						r.Check(tr >= 1,
							fmt.Sprintf("case-insensitive filter matched %v results, want >= 1", tr))
					}
				},
			},
		},
	},

	// Extensions.
	{
		ID:       "EXTRA-EXT-CONTAINER",
		Level:    Must,
		Summary:  "Extension attributes must be namespaced under the extension schema URI",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "extension_container",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"schemas": []string{TestResourceSchema, TestExtSchema},
						TestExtSchema: map[string]any{
							"extString":   "ext-value",
							"extRequired": "required-value",
						},
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					extData, ok := body[TestExtSchema].(map[string]any)
					r.Check(ok && extData != nil,
						fmt.Sprintf("extension data not found under key %q", TestExtSchema))
				},
			},
		},
	},
	{
		ID:       "EXTRA-EXT-ROUNDTRIP",
		Level:    Must,
		Summary:  "Extension attributes must round-trip through create and get",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "extension_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"schemas": []string{TestResourceSchema, TestExtSchema},
						TestExtSchema: map[string]any{
							"extString":   "ext-value",
							"extRequired": "required-value",
						},
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					extData, ok := body[TestExtSchema].(map[string]any)
					if ok && extData != nil {
						es, _ := extData["extString"].(string)
						r.Check(es == "ext-value",
							fmt.Sprintf("extString = %q, want \"ext-value\"", es))

						er, _ := extData["extRequired"].(string)
						r.Check(er == "required-value",
							fmt.Sprintf("extRequired = %q, want \"required-value\"", er))
					}

					id := IDOf(body)
					getResp, err := r.Client.Get("/TestResources/" + id)
					r.RequireOK(err)

					if getResp.StatusCode == 200 {
						extGet, ok := getResp.Body[TestExtSchema].(map[string]any)
						r.Check(ok && extGet != nil,
							"GET: extension data missing after creation")
					}
				},
			},
		},
	},
	{
		ID:       "EXTRA-EXT-SCHEMAS",
		Level:    Must,
		Summary:  "Resources with extensions must include extension URI in schemas array",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "extension_schemas",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"schemas": []string{TestResourceSchema, TestExtSchema},
						TestExtSchema: map[string]any{
							"extString":   "ext-value",
							"extRequired": "required-value",
						},
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					schemas := GetStringSlice(body, "schemas")
					r.Check(HasSchema(body, TestExtSchema),
						fmt.Sprintf("schemas = %v, missing extension URI %s", schemas, TestExtSchema))
				},
			},
		},
	},

	// Multi-valued.
	{
		ID:       "EXTRA-MV-ROUNDTRIP",
		Level:    Must,
		Summary:  "Multi-valued attributes must round-trip preserving all values",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "multi_valued_round_trip",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"multiStrings": []string{"alpha", "beta", "gamma"},
					})
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}

					arr, ok := body["multiStrings"].([]any)
					r.Check(ok && len(arr) == 3,
						fmt.Sprintf("multiStrings = %v, want 3 elements", body["multiStrings"]))
				},
			},
		},
	},
	{
		ID:       "EXTRA-MV-COMPLEX-PRIMARY",
		Level:    Must,
		Summary:  "Multi-valued complex attributes must enforce at most one primary=true",
		Feature:  TestResource,
		Testable: true,
		Tests: []Test{
			{
				Name: "multi_complex_primary",
				Fn: func(r *Run) {
					body, resp := r.CreateTestResource(map[string]any{
						"multiComplex": []map[string]any{
							{"value": "one", "type": "a", "primary": true},
							{"value": "two", "type": "b", "primary": true},
						},
					})
					if resp.StatusCode == 201 {
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
						r.Check(primaryCount <= 1,
							fmt.Sprintf("multiComplex has %d primary=true, want <= 1", primaryCount))
					} else if resp.StatusCode == 400 {
						// Rejecting duplicate primary is also acceptable.
						r.Check(true, "")
					} else {
						r.Check(false,
							FmtStatus("POST /TestResources with duplicate primary", resp.StatusCode, 400))
					}
				},
			},
		},
	},
}

func init() {
	All = append(All, Extra...)
}
