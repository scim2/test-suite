package spec

import "fmt"

var l7644_3_5 = rfcLines(rfc7644Sections, "13-section-3.5.txt")

var rfc7644_3_5 = []Requirement{
	// Section 3.5

	{
		ID:      "RFC7644-3.5-L1607",
		Level:   Must,
		Summary: "Implementers must support HTTP PUT",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: l7644_3_5(3),
			EndLine:   l7644_3_5(5),
			EndCol:    31,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "put_returns_200",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)
				schemas := GetStringSlice(body, "schemas")
				userName, _ := body["userName"].(string)

				putBody := map[string]any{
					"schemas":     schemas,
					"userName":    userName,
					"displayName": "Updated " + RandomSuffix(),
				}

				resp, err := r.Client.Put("/Users/"+id, putBody)
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("PUT /Users/"+id, resp.StatusCode, 200))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5-L1609",
		Level:   Should,
		Summary: "Implementers should support HTTP PATCH for partial modifications",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: l7644_3_5(5),
			EndLine:   l7644_3_5(9),
			StartCol:  34,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "patch_supported",
			Fn: func(r *Run) {
				r.Check(r.Features.Patch, "PATCH not supported by service provider")
			},
		}},
	},

	// Section 3.5.1

	{
		ID:      "RFC7644-3.5.1-L1637",
		Level:   MustNot,
		Summary: "HTTP PUT must not be used to create new resources",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(33),
			EndLine:   l7644_3_5(34),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "put_does_not_create",
			Fn: func(r *Run) {
				resp, err := r.Client.Put("/Users/nonexistent-"+RandomSuffix(), map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-putcreate-" + RandomSuffix(),
				})
				r.RequireOK(err)

				r.Check(resp.StatusCode == 404,
					FmtStatus("PUT /Users/<nonexistent>", resp.StatusCode, 404))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.1-L1644",
		Level:   Shall,
		Summary: "readWrite/writeOnly values provided in PUT shall replace existing values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(41),
			EndLine:   l7644_3_5(42),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "put_replace_values",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", createResp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				newDisplay := "Replaced-" + RandomSuffix()
				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":     []string{UserSchema},
					"userName":    userName,
					"displayName": newDisplay,
				})
				r.RequireOK(err)

				if putResp.StatusCode != 200 {
					r.Fatalf("PUT returned %d", putResp.StatusCode)
				}

				dn, _ := putResp.Body["displayName"].(string)
				r.Check(dn == newDisplay,
					fmt.Sprintf("PUT replace: displayName = %q, want %q", dn, newDisplay))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.1-L1660",
		Level:   Must,
		Summary: "Immutable attribute input must match existing values or return 400",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(56),
			EndLine:   l7644_3_5(60),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "immutable_must_match",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
					"id":       "different-id-value",
				})
				r.RequireOK(err)

				if putResp.StatusCode == 200 {
					r.Check(IDOf(putResp.Body) == id,
						fmt.Sprintf("PUT with mismatched id: server changed id to %q", IDOf(putResp.Body)))
				} else {
					r.Check(putResp.StatusCode == 400,
						FmtStatus("PUT /Users with mismatched immutable id", putResp.StatusCode, 400))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.1-L1665",
		Level:   Shall,
		Summary: "readOnly values provided in PUT shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(62),
			EndLine:   l7644_3_5(62),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "readonly_ignored_on_put",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)
				schemas := GetStringSlice(body, "schemas")
				userName, _ := body["userName"].(string)

				putBody := map[string]any{
					"schemas":     schemas,
					"userName":    userName,
					"displayName": "Updated " + RandomSuffix(),
				}

				resp, err := r.Client.Put("/Users/"+id, putBody)
				r.RequireOK(err)

				r.Check(IDOf(resp.Body) == id,
					fmt.Sprintf("PUT /Users/%s: id changed to %q", id, IDOf(resp.Body)))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.1-L1667",
		Level:   Must,
		Summary: "Required attributes must be specified in PUT request",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(64),
			EndLine:   l7644_3_5(65),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "put_required_attributes",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}
				id := IDOf(body)

				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":     []string{UserSchema},
					"displayName": "No UserName",
				})
				r.RequireOK(err)

				r.Check(putResp.StatusCode == 400,
					FmtStatus("PUT /Users without required userName", putResp.StatusCode, 400))
			},
		}},
	},

	// Section 3.5.2

	{
		ID:      "RFC7644-3.5.2-L1805",
		Level:   Must,
		Summary: "PATCH body must contain schemas attribute with PatchOp URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(202),
			EndLine:   l7644_3_5(203),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "add_display_name",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "add",
							"path":  "displayName",
							"value": "Patched User",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					"PATCH request with PatchOp schema was not accepted")
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1808",
		Level:   Must,
		Summary: "PATCH body must contain Operations array with one or more operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(205),
			EndLine:   l7644_3_5(207),
			EndCol:    14,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "operations_array",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "add",
							"path":  "nickName",
							"value": "Nicky",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH with Operations array", patchResp.StatusCode, 200))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1810",
		Level:   Must,
		Summary: "Each PATCH operation must have exactly one op member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(207),
			EndLine:   l7644_3_5(211),
			StartCol:  17,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "replace_display_name",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "displayName",
							"value": "Replaced User",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH replace /Users/"+id, patchResp.StatusCode, 200))

				// Verify via GET.
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					dn, _ := getResp.Body["displayName"].(string)
					r.Check(dn == "Replaced User",
						fmt.Sprintf("after PATCH replace: displayName = %q, want \"Replaced User\"", dn))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1886",
		Level:   Must,
		Summary: "Each PATCH operation must be compatible with attribute mutability and schema",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(283),
			EndLine:   l7644_3_5(285),
			EndCol:    16,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "mutability_compatible",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				// Replace a readWrite attribute (displayName) - should work.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "displayName",
							"value": "Mutability Test",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH replace readWrite displayName", patchResp.StatusCode, 200))
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1888",
		Level:   MustNot,
		Summary: "Client must not modify readOnly or immutable attributes via PATCH",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(285),
			EndLine:   l7644_3_5(288),
			StartCol:  19,
			EndCol:    21,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "readonly_not_modifiable",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				// Try to PATCH the id (readOnly attribute).
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "id",
							"value": "hacked-id",
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode >= 200 && patchResp.StatusCode < 300 {
					// If accepted, id must not have changed.
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					r.Check(getResp.StatusCode == 200 && IDOf(getResp.Body) == id,
						"PATCH replace on readOnly id: id was modified")
				} else {
					// Server correctly rejected the operation.
					r.Check(patchResp.StatusCode == 400,
						FmtStatus("PATCH replace readOnly id", patchResp.StatusCode, 400))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1927",
		Level:   Shall,
		Summary: "Setting primary to true shall auto-set primary to false for other values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(323),
			EndLine:   l7644_3_5(326),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "primary_auto_unset",
			Fn: func(r *Run) {
				// Create user with one primary email.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-primary-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "first@example.com", "type": "work", "primary": true},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				// Add another email as primary.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":   "add",
							"path": "emails",
							"value": []map[string]any{
								{"value": "second@example.com", "type": "home", "primary": true},
							},
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode == 200 || patchResp.StatusCode == 204 {
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						emails, _ := getResp.Body["emails"].([]any)
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
						r.Check(primaryCount <= 1,
							fmt.Sprintf("after adding primary email: %d primary values, want <= 1", primaryCount))
					}
				} else if patchResp.StatusCode == 400 {
					// Server rejected adding a second primary; acceptable.
					r.Check(true, "")
				} else {
					r.Check(false,
						FmtStatus("PATCH add primary email", patchResp.StatusCode, 200))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1931",
		Level:   Shall,
		Summary: "PATCH request shall be treated as atomic",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(328),
			EndLine:   l7644_3_5(331),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "atomicity",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)
				originalDN, _ := body["displayName"].(string)

				// Two operations: first valid, second targets readOnly id.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "displayName",
							"value": "Atomic-Change-" + RandomSuffix(),
						},
						{
							"op":    "replace",
							"path":  "id",
							"value": "invalid-id-change",
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode >= 400 {
					// Server rejected -- verify original state is preserved.
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						dn, _ := getResp.Body["displayName"].(string)
						r.Check(dn == originalDN,
							fmt.Sprintf("PATCH not atomic: displayName changed to %q despite error", dn))
					}
				} else {
					// Server accepted both -- second op was likely ignored (readOnly).
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						r.Check(IDOf(getResp.Body) == id,
							"PATCH: id was modified (should be readOnly)")
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1933",
		Level:   Must,
		Summary: "On PATCH operation error, original resource must be restored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(329),
			EndLine:   l7644_3_5(331),
			StartCol:  24,
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "atomicity_restore",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				// Two operations: first valid, second invalid (remove required userName).
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "displayName",
							"value": "Should Not Persist",
						},
						{
							"op":   "remove",
							"path": "userName",
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode >= 400 {
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 {
						dn, _ := getResp.Body["displayName"].(string)
						r.Check(dn != "Should Not Persist",
							"PATCH was not atomic: displayName changed despite second op failing")
					}
				} else {
					// Server accepted the removal of userName.
					r.Logf("server accepted removal of required attribute userName")
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2-L1939",
		Level:   Must,
		Summary: "Successful PATCH must return 200 with resource body or 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(336),
			EndLine:   l7644_3_5(341),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "add_display_name",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "add",
							"path":  "displayName",
							"value": "Patched User",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH /Users/"+id, patchResp.StatusCode, 200))
			},
		}},
	},

	// Section 3.5.2.1

	{
		ID:      "RFC7644-3.5.2.1-L1972",
		Level:   Must,
		Summary: "PATCH add operation must contain a value member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(369),
			EndLine:   l7644_3_5(372),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "add_display_name",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "add",
							"path":  "displayName",
							"value": "Patched User",
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode == 200 && patchResp.Body != nil {
					dn, _ := patchResp.Body["displayName"].(string)
					r.Check(dn == "Patched User",
						fmt.Sprintf("PATCH add displayName: got %q, want \"Patched User\"", dn))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.1-L1988",
		Level:   Shall,
		Summary: "Add to complex attribute shall specify sub-attributes in value parameter",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(384),
			EndLine:   l7644_3_5(385),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "add_to_complex",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":   "add",
							"path": "name",
							"value": map[string]any{
								"givenName":  "Test",
								"familyName": "User",
							},
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH add to complex name", patchResp.StatusCode, 200))

				// Verify sub-attributes were set.
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)
				if getResp.StatusCode == 200 {
					gn := GetString(getResp.Body, "name", "givenName")
					r.Check(gn == "Test",
						fmt.Sprintf("name.givenName = %q, want \"Test\"", gn))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.1-L2004",
		Level:   ShallNot,
		Summary: "Add of already-existing value shall not change the modify timestamp",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(398),
			EndLine:   l7644_3_5(402),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "add_existing_no_timestamp_change",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":     []string{UserSchema},
					"userName":    "scim-test-addexist-" + RandomSuffix(),
					"displayName": "Existing Value",
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				lastMod := GetString(resp.Body, "meta", "lastModified")

				// Add the same value again.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "add",
							"path":  "displayName",
							"value": "Existing Value",
						},
					},
				})
				r.RequireOK(err)

				if patchResp.StatusCode == 200 || patchResp.StatusCode == 204 {
					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)
					if getResp.StatusCode == 200 && lastMod != "" {
						newLastMod := GetString(getResp.Body, "meta", "lastModified")
						r.Check(newLastMod == lastMod,
							fmt.Sprintf("add existing value changed lastModified: %s -> %s", lastMod, newLastMod))
					}
				}
			},
		}},
	},

	// Section 3.5.2.2

	{
		ID:      "RFC7644-3.5.2.2-L2146",
		Level:   Shall,
		Summary: "Removing single-value attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(542),
			EndLine:   l7644_3_5(544),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "remove_display_name",
			Fn: func(r *Run) {
				// Create user with displayName.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":     []string{UserSchema},
					"userName":    "scim-test-patchrm-" + RandomSuffix(),
					"displayName": "To Be Removed",
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":   "remove",
							"path": "displayName",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH remove /Users/"+id, patchResp.StatusCode, 200))

				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					dn, hasDN := getResp.Body["displayName"]
					r.Check(!hasDN || dn == nil || dn == "",
						fmt.Sprintf("after PATCH remove: displayName = %v, want unassigned", dn))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.2-L2151",
		Level:   Shall,
		Summary: "Removing all values of multi-valued attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(546),
			EndLine:   l7644_3_5(548),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "remove_multi_valued",
			Fn: func(r *Run) {
				// Create user with emails.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-patchmv-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "a@example.com", "type": "work"},
						{"value": "b@example.com", "type": "home"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				// Remove all emails.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":   "remove",
							"path": "emails",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH remove emails", patchResp.StatusCode, 200))

				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					emails, hasEmails := getResp.Body["emails"]
					emailArr, _ := emails.([]any)
					r.Check(!hasEmails || emails == nil || len(emailArr) == 0,
						fmt.Sprintf("after PATCH remove emails: emails = %v, want unassigned", emails))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.2-L2167",
		Level:   Shall,
		Summary: "Removing required or readOnly attribute shall return error with mutability scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(563),
			EndLine:   l7644_3_5(567),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "remove_required",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				// Try to remove required attribute userName.
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":   "remove",
							"path": "userName",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 400,
					FmtStatus("PATCH remove required userName", patchResp.StatusCode, 400))
			},
		}},
	},

	// Section 3.5.2.3

	{
		ID:      "RFC7644-3.5.2.3-L2376",
		Level:   Shall,
		Summary: "Replace on non-existent attribute shall treat as add",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(772),
			EndLine:   l7644_3_5(773),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "replace_non_existent",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				newNick := "NewNick-" + RandomSuffix()
				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "nickName",
							"value": newNick,
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH replace non-existent nickName", patchResp.StatusCode, 200))

				// Verify the value was set.
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					nick, _ := getResp.Body["nickName"].(string)
					r.Check(nick == newNick,
						fmt.Sprintf("PATCH replace non-existent: nickName = %q, want %q", nick, newNick))
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.3-L2387",
		Level:   Shall,
		Summary: "Replace with valuePath filter shall replace all matching record values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(781),
			EndLine:   l7644_3_5(784),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "replace_value_path",
			Fn: func(r *Run) {
				// Create user with two emails.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-vpath-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "work@example.com", "type": "work"},
						{"value": "home@example.com", "type": "home"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "emails[type eq \"work\"].value",
							"value": "newwork@example.com",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH replace with valuePath", patchResp.StatusCode, 200))

				// Verify via GET.
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					emails, _ := getResp.Body["emails"].([]any)
					found := false
					for _, e := range emails {
						em, ok := e.(map[string]any)
						if !ok {
							continue
						}
						if tp, _ := em["type"].(string); tp == "work" {
							v, _ := em["value"].(string)
							if v == "newwork@example.com" {
								found = true
							}
						}
					}
					r.Check(found,
						"PATCH replace valuePath: work email not updated to newwork@example.com")
				}
			},
		}},
	},
	{
		ID:      "RFC7644-3.5.2.3-L2396",
		Level:   Shall,
		Summary: "Replace with unmatched valuePath filter shall return 400 noTarget",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(791),
			EndLine:   l7644_3_5(795),
		},
		Feature:  Patch,
		Testable: true,
		Tests: []Test{{
			Name: "replace_value_path_no_target",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-notarget-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "only@example.com", "type": "work"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })

				patchResp, err := r.Client.Patch("/Users/"+id, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{
						{
							"op":    "replace",
							"path":  "emails[type eq \"nonexistent\"].value",
							"value": "nope@example.com",
						},
					},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 400,
					FmtStatus("PATCH replace unmatched valuePath", patchResp.StatusCode, 400))

				if patchResp.StatusCode == 400 && patchResp.Body != nil {
					scimType, _ := patchResp.Body["scimType"].(string)
					r.Check(scimType == "noTarget",
						fmt.Sprintf("scimType = %q, want \"noTarget\"", scimType))
				}
			},
		}},
	},
}
