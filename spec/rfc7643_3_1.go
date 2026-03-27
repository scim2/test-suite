package spec

import (
	"fmt"
	"strings"
)

var l7643_3_1 = rfcLines(rfc7643Sections, "11-section-3.1.txt")

var rfc7643_3_1 = []Requirement{
	// Section 3.1

	{
		ID:      "RFC7643-3.1-L852",
		Level:   Must,
		Summary: "Common attributes must be defined for all resources (except ServiceProviderConfig and ResourceType)",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(2),
			EndLine:   l7643_3_1(7),
			EndCol:    41,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "meta_attribute_present",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				meta, _ := body["meta"].(map[string]any)
				r.Check(meta != nil,
					"POST /Users: missing meta attribute")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L855",
		Level:   Must,
		Summary: "Service provider must assign id and meta attributes upon creation",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(7),
			EndLine:   l7643_3_1(11),
			StartCol:  44,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "id_and_meta_assigned",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				r.Check(IDOf(body) != "" && body["meta"] != nil,
					"POST /Users: service provider must assign id and meta on creation")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L866",
		Level:   Must,
		Summary: "Each resource representation must include a non-empty id value",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(19),
			EndLine:   l7643_3_1(21),
			EndCol:    27,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "id_non_empty",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				r.Check(IDOf(body) != "",
					"POST /Users: missing id in response body")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L867",
		Level:   Must,
		Summary: "Resource id must be unique across the entire set of resources",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(21),
			EndLine:   l7643_3_1(22),
			StartCol:  30,
			EndCol:    54,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "id_unique",
			Fn: func(r *Run) {
				_, resp1 := r.CreateUser()
				_, resp2 := r.CreateUser()

				if resp1.StatusCode != 201 || resp2.StatusCode != 201 {
					r.Skipf("setup: could not create two users")
				}

				r.Check(IDOf(resp1.Body) != IDOf(resp2.Body),
					fmt.Sprintf("two users got same id: %q", IDOf(resp1.Body)))
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L868",
		Level:   Must,
		Summary: "Resource id must be stable and non-reassignable",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(22),
			EndLine:   l7643_3_1(24),
			StartCol:  57,
			EndCol:    55,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "id_stable_on_get",
			Fn: func(r *Run) {
				body, createResp := r.CreateUser()
				if createResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
				}
				id := IDOf(body)

				resp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				r.Check(IDOf(resp.Body) == id,
					fmt.Sprintf("GET /Users/%s: id = %q, want %q", id, IDOf(resp.Body), id))
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L872",
		Level:   MustNot,
		Summary: "Resource id must not be specified by the client",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(24),
			EndLine:   l7643_3_1(26),
			StartCol:  57,
			EndCol:    42,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "client_id_ignored",
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
					"POST /Users: server accepted client-specified id")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L873",
		Level:   MustNot,
		Summary: "The string bulkId must not be used within any unique identifier value",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(26),
			EndLine:   l7643_3_1(28),
			StartCol:  45,
			EndCol:    12,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "id_no_bulkId",
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

				id := IDOf(resp.Body)
				r.Check(id != "bulkId" && !strings.Contains(id, "bulkId"),
					fmt.Sprintf("POST /Users: id %q contains reserved keyword bulkId", id))
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L890",
		Level:   MustNot,
		Summary: "externalId must not be specified by the service provider",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(42),
			EndLine:   l7643_3_1(44),
			StartCol:  47,
			EndCol:    56,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "external_id_not_set",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				_, hasExtId := body["externalId"]
				r.Check(!hasExtId,
					"POST /Users without externalId: server added externalId to response")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L891",
		Level:   Must,
		Summary: "Service provider must interpret externalId as scoped to the provisioning domain",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(44),
			EndLine:   l7643_3_1(46),
			StartCol:  59,
			EndCol:    26,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-3.1-L911",
		Level:   Shall,
		Summary: "The meta attribute shall be ignored when provided by clients",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(64),
			EndLine:   l7643_3_1(66),
			EndCol:    39,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "client_meta_ignored",
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

				rt := GetString(resp.Body, "meta", "resourceType")
				r.Check(rt == "User",
					fmt.Sprintf("POST /Users: meta.resourceType = %q, server should have ignored client meta", rt))
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L920",
		Level:   Must,
		Summary: "meta.created must be a DateTime",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(72),
			EndLine:   l7643_3_1(74),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "meta_created_present",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				created := GetString(body, "meta", "created")
				r.Check(created != "",
					"POST /Users: missing meta.created")
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L925",
		Level:   Must,
		Summary: "meta.lastModified must equal meta.created if resource was never modified",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(75),
			EndLine:   l7643_3_1(79),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "created_equals_lastModified",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				created := GetString(body, "meta", "created")
				lastMod := GetString(body, "meta", "lastModified")
				r.Check(created == lastMod,
					fmt.Sprintf("POST /Users: meta.created (%s) != meta.lastModified (%s) on new resource",
						created, lastMod))
			},
		}},
	},
	{
		ID:      "RFC7643-3.1-L927",
		Level:   Must,
		Summary: "meta.location must match Content-Location header",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(80),
			EndLine:   l7643_3_1(83),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{
			{
				Name: "create_content_location_matches_meta",
				Fn: func(r *Run) {
					body, resp := r.CreateUser()
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
					}

					meta, _ := body["meta"].(map[string]any)
					if meta != nil {
						contentLoc := resp.Header.Get("Content-Location")
						metaLoc, _ := meta["location"].(string)
						if contentLoc != "" {
							r.Check(contentLoc == metaLoc,
								fmt.Sprintf("Content-Location (%s) != meta.location (%s)", contentLoc, metaLoc))
						}
					}
				},
			},
			{
				Name: "get_content_location_matches_meta",
				Fn: func(r *Run) {
					body, createResp := r.CreateUser()
					if createResp.StatusCode != 201 {
						r.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
					}
					id := IDOf(body)

					resp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)

					contentLoc := resp.Header.Get("Content-Location")
					metaLoc := GetString(resp.Body, "meta", "location")
					if contentLoc != "" && metaLoc != "" {
						r.Check(contentLoc == metaLoc,
							fmt.Sprintf("Content-Location (%s) != meta.location (%s)",
								contentLoc, metaLoc))
					}
				},
			},
			{
				Name: "meta_location_present",
				Fn: func(r *Run) {
					body, resp := r.CreateUser()
					if resp.StatusCode != 201 {
						r.Fatalf("setup: POST returned %d", resp.StatusCode)
					}
					id := IDOf(body)

					getResp, err := r.Client.Get("/Users/" + id)
					r.RequireOK(err)

					if getResp.StatusCode != 200 {
						r.Fatalf("GET /Users/%s returned %d", id, getResp.StatusCode)
					}

					metaLoc := GetString(getResp.Body, "meta", "location")
					contentLoc := getResp.Header.Get("Content-Location")
					r.Check(metaLoc != "",
						"GET /Users: missing meta.location")

					if contentLoc != "" && metaLoc != "" {
						r.Check(contentLoc == metaLoc,
							fmt.Sprintf("Content-Location (%s) != meta.location (%s)", contentLoc, metaLoc))
					}
				},
			},
		},
	},
	{
		ID:      "RFC7643-3.1-L940",
		Level:   Must,
		Summary: "Weak entity-tags must be prefixed with W/",
		Source: Source{
			RFC:       7643,
			Section:   "3.1",
			StartLine: l7643_3_1(89),
			EndLine:   l7643_3_1(96),
		},
		Feature:  ETag,
		Testable: true,
		Tests: []Test{{
			Name: "weak_etag_prefix",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				version := GetString(body, "meta", "version")
				if version != "" {
					r.Check(strings.HasPrefix(version, "W/"),
						fmt.Sprintf("meta.version = %q, weak ETags must start with W/", version))
				}

				id := IDOf(body)
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				etag := getResp.Header.Get("ETag")
				if etag != "" {
					r.Check(strings.HasPrefix(etag, "W/"),
						fmt.Sprintf("ETag header = %q, weak ETags must start with W/", etag))
				}
			},
		}},
	},
}
