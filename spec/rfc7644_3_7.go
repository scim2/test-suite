package spec

import "fmt"

var l7644_3_7 = rfcLines(rfc7644Sections, "15-section-3.7.txt")

var rfc7644_3_7 = []Requirement{
	// Section 3.7

	{
		ID:      "RFC7644-3.7-L2775",
		Level:   Must,
		Summary: "Successful bulk job must return 200 OK",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(81),
			EndLine:   l7644_3_7(83),
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "returns_200",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas": []string{BulkSchema},
					"Operations": []map[string]any{{
						"method": "POST",
						"path":   "/Users",
						"bulkId": "test-1",
						"data": map[string]any{
							"schemas":  []string{UserSchema},
							"userName": "scim-test-bulk-" + RandomSuffix(),
						},
					}},
				})
				r.RequireOK(err)
				r.Check(resp.StatusCode == 200, FmtStatus("POST /Bulk", resp.StatusCode, 200))

				if resp.StatusCode == 200 && resp.Body != nil {
					cleanupBulkOps(r, resp.Body)
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.7-L2779",
		Level:   Must,
		Summary: "Service provider must continue performing changes and disregard partial failures",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(85),
			EndLine:   l7644_3_7(86),
			EndCol:    43,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "continues_on_partial_failure",
			Fn: func(r *Run) {
				userName := "scim-test-bulkdup-" + RandomSuffix()

				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas":      []string{BulkSchema},
					"failOnErrors": 10,
					"Operations": []map[string]any{
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "ok-1",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": userName,
							},
						},
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "dup-1",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": userName,
							},
						},
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "ok-2",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": "scim-test-bulkok-" + RandomSuffix(),
							},
						},
					},
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("POST /Bulk (partial failure)", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					ops, _ := resp.Body["Operations"].([]any)
					r.Check(len(ops) == 3,
						fmt.Sprintf("bulk returned %d operations, want 3 (all processed)", len(ops)))
					cleanupBulkOps(r, resp.Body)
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.7-L2789",
		Level:   Must,
		Summary: "Service provider must return the same bulkId with newly created resource",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(95),
			EndLine:   l7644_3_7(98),
			StartCol:  36,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "returns_bulk_id",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas": []string{BulkSchema},
					"Operations": []map[string]any{{
						"method": "POST",
						"path":   "/Users",
						"bulkId": "myBulkId-42",
						"data": map[string]any{
							"schemas":  []string{UserSchema},
							"userName": "scim-test-bulkid-" + RandomSuffix(),
						},
					}},
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("POST /Bulk (bulkId)", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					ops, _ := resp.Body["Operations"].([]any)
					if len(ops) > 0 {
						first, _ := ops[0].(map[string]any)
						bulkId, _ := first["bulkId"].(string)
						r.Check(bulkId == "myBulkId-42",
							fmt.Sprintf("bulkId = %q, want \"myBulkId-42\"", bulkId))
					}
					cleanupBulkOps(r, resp.Body)
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.7-L2796",
		Level:   Must,
		Summary: "Optimized bulk processing must preserve client intent and achieve same result",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(101),
			EndLine:   l7644_3_7(113),
			StartCol:  68,
			EndCol:    31,
		},
		Feature:  Bulk,
		Testable: false,
	},

	// Section 3.7.1

	{
		ID:      "RFC7644-3.7.1-L2814",
		Level:   Must,
		Summary: "Service provider must try to resolve circular cross-references in bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.1",
			StartLine: l7644_3_7(120),
			EndLine:   l7644_3_7(122),
			EndCol:    62,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "circular_references",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas": []string{BulkSchema},
					"Operations": []map[string]any{
						{
							"method": "POST",
							"path":   "/Groups",
							"bulkId": "groupA",
							"data": map[string]any{
								"schemas":     []string{GroupSchema},
								"displayName": "scim-test-circA-" + RandomSuffix(),
								"members": []map[string]any{
									{"value": "bulkId:groupB", "type": "Group"},
								},
							},
						},
						{
							"method": "POST",
							"path":   "/Groups",
							"bulkId": "groupB",
							"data": map[string]any{
								"schemas":     []string{GroupSchema},
								"displayName": "scim-test-circB-" + RandomSuffix(),
								"members": []map[string]any{
									{"value": "bulkId:groupA", "type": "Group"},
								},
							},
						},
					},
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("POST /Bulk with circular refs", resp.StatusCode, 200))

				if resp.StatusCode == 200 && resp.Body != nil {
					cleanupBulkOps(r, resp.Body)
				}
			},
		}},
	},

	// Section 3.7.3

	{
		ID:      "RFC7644-3.7.3-L3201",
		Level:   Must,
		Summary: "Bulk response must include result of all processed operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: l7644_3_7(507),
			EndLine:   l7644_3_7(512),
			EndCol:    33,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "response_includes_all_ops",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas": []string{BulkSchema},
					"Operations": []map[string]any{
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "r1",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": "scim-test-bulkall1-" + RandomSuffix(),
							},
						},
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "r2",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": "scim-test-bulkall2-" + RandomSuffix(),
							},
						},
					},
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("POST /Bulk (response ops)", resp.StatusCode, 200))

				if resp.StatusCode == 200 {
					ops, _ := resp.Body["Operations"].([]any)
					r.Check(len(ops) == 2,
						fmt.Sprintf("bulk response has %d operations, want 2", len(ops)))
					cleanupBulkOps(r, resp.Body)
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.7.3-L3206",
		Level:   Must,
		Summary: "Bulk status must include HTTP response code",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: l7644_3_7(511),
			EndLine:   l7644_3_7(516),
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "status_code_present",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas": []string{BulkSchema},
					"Operations": []map[string]any{
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "s1",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": "scim-test-bulkst1-" + RandomSuffix(),
							},
						},
						{
							"method": "POST",
							"path":   "/Users",
							"bulkId": "s2",
							"data": map[string]any{
								"schemas":  []string{UserSchema},
								"userName": "scim-test-bulkst2-" + RandomSuffix(),
							},
						},
					},
				})
				r.RequireOK(err)

				if resp.StatusCode != 200 {
					r.Check(false,
						FmtStatus("POST /Bulk (status codes)", resp.StatusCode, 200))
					return
				}

				ops, _ := resp.Body["Operations"].([]any)
				for i, op := range ops {
					o, ok := op.(map[string]any)
					if !ok {
						continue
					}
					status, _ := o["status"].(string)
					r.Check(status != "",
						fmt.Sprintf("bulk operation %d missing status", i))
				}

				cleanupBulkOps(r, resp.Body)
			},
		}},
	},

	// Section 3.7.4

	{
		ID:      "RFC7644-3.7.4-L3495",
		Level:   Must,
		Summary: "Service provider must define max operations and max payload size for bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: l7644_3_7(801),
			EndLine:   l7644_3_7(804),
			EndCol:    50,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "max_limits_defined",
			Fn: func(r *Run) {
				spcResp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				bulkCfg, _ := spcResp.Body["bulk"].(map[string]any)
				maxOps, hasMaxOps := bulkCfg["maxOperations"]
				maxPayload, hasMaxPayload := bulkCfg["maxPayloadSize"]

				r.Check(hasMaxOps && hasMaxPayload,
					fmt.Sprintf("bulk config missing limits: maxOperations=%v, maxPayloadSize=%v", maxOps, maxPayload))
			},
		}},
	},
	{
		ID:      "RFC7644-3.7.4-L3499",
		Level:   Must,
		Summary: "Exceeding bulk limits must return 413 Payload Too Large",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: l7644_3_7(804),
			EndLine:   l7644_3_7(807),
			StartCol:  49,
		},
		Feature:  Bulk,
		Testable: true,
		Tests: []Test{{
			Name: "exceed_limits_returns_413",
			Fn: func(r *Run) {
				spcResp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				bulkCfg, _ := spcResp.Body["bulk"].(map[string]any)
				maxOps, _ := bulkCfg["maxOperations"].(float64)
				if maxOps == 0 {
					maxOps = 10
				}

				ops := make([]map[string]any, int(maxOps)+1)
				for i := range ops {
					ops[i] = map[string]any{
						"method": "POST",
						"path":   "/Users",
						"bulkId": fmt.Sprintf("excess-%d", i),
						"data": map[string]any{
							"schemas":  []string{UserSchema},
							"userName": fmt.Sprintf("scim-test-exceed-%d-%s", i, RandomSuffix()),
						},
					}
				}

				resp, err := r.Client.Post("/Bulk", map[string]any{
					"schemas":    []string{BulkSchema},
					"Operations": ops,
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 413,
					FmtStatus("POST /Bulk exceeding maxOperations", resp.StatusCode, 413))
			},
		}},
	},
}

// cleanupBulkOps deletes any resources created by a bulk response.
func cleanupBulkOps(r *Run, body map[string]any) {
	ops, _ := body["Operations"].([]any)
	for _, op := range ops {
		o, ok := op.(map[string]any)
		if !ok {
			continue
		}
		if loc, _ := o["location"].(string); loc != "" {
			_, _ = r.Client.Delete(loc)
		}
	}
}
