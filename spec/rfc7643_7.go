package spec

import (
	"fmt"
	"net/url"
)

var l7643_7 = rfcLines(rfc7643Sections, "20-section-7.txt")

var rfc7643_7 = []Requirement{
	// Section 7

	{
		ID:      "RFC7643-7-L1691",
		Level:   Must,
		Summary: "Service providers must specify the schema URI",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(31),
			EndLine:   l7643_7(33),
			EndCol:    51,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "schema_has_id",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Schemas")
				r.RequireOK(err)

				if resp.StatusCode != 200 || resp.Body == nil {
					r.Fatalf("GET /Schemas returned %d", resp.StatusCode)
				}

				resources, _ := resp.Body["Resources"].([]any)
				for _, res := range resources {
					s, ok := res.(map[string]any)
					if !ok {
						continue
					}
					id, _ := s["id"].(string)
					r.Check(id != "",
						"Schema missing id (URI)")
				}
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1700",
		Level:   Must,
		Summary: "Service providers must specify the schema name when applicable",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(40),
			EndLine:   l7643_7(41),
			StartCol:  42,
			EndCol:    63,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "schema_has_name",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Schemas")
				r.RequireOK(err)

				if resp.StatusCode != 200 || resp.Body == nil {
					r.Fatalf("GET /Schemas returned %d", resp.StatusCode)
				}

				resources, _ := resp.Body["Resources"].([]any)
				for _, res := range resources {
					s, ok := res.(map[string]any)
					if !ok {
						continue
					}
					id, _ := s["id"].(string)
					name, _ := s["name"].(string)
					r.Check(name != "",
						fmt.Sprintf("Schema %q: missing name", id))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1748",
		Level:   Shall,
		Summary: "Server shall use case sensitivity when evaluating filters per caseExact",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(88),
			EndLine:   l7643_7(90),
			EndCol:    45,
		},
		Feature:  Filter,
		Testable: true,
		Tests: []Test{{
			Name: "filter_case_exact",
			Fn: func(r *Run) {
				_, resp := r.CreateTestResource(map[string]any{
					"caseExactString": "CaSeExAcT",
				})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				// Search with wrong case should NOT match.
				filter := "caseExactString eq \"caseexact\""
				qResp, err := r.Client.Get("/TestResources?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				if qResp.StatusCode == 200 {
					tr, _ := qResp.Body["totalResults"].(float64)
					r.Check(tr == 0,
						fmt.Sprintf("caseExact filter with wrong case matched %v results, want 0", tr))
				}

				// Search with correct case should match.
				filter = "caseExactString eq \"CaSeExAcT\""
				qResp2, err := r.Client.Get("/TestResources?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				if qResp2.StatusCode == 200 {
					tr, _ := qResp2.Body["totalResults"].(float64)
					r.Check(tr >= 1,
						fmt.Sprintf("caseExact filter with correct case matched %v results, want >= 1", tr))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1750",
		Level:   Shall,
		Summary: "Server shall preserve case for case-exact attribute values",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(90),
			EndLine:   l7643_7(92),
			StartCol:  48,
			EndCol:    19,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "case_preservation",
			Fn: func(r *Run) {
				mixedCase := "ScIm-TeSt-CaSe-" + RandomSuffix()
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":     []string{UserSchema},
					"userName":    mixedCase,
					"displayName": "CaSeExAcT Display",
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode != 201 {
					r.Skipf("POST returned %d", resp.StatusCode)
				}

				dn, _ := resp.Body["displayName"].(string)
				r.Check(dn == "CaSeExAcT Display",
					fmt.Sprintf("displayName = %q, want \"CaSeExAcT Display\" (case not preserved)", dn))
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1759",
		Level:   ShallNot,
		Summary: "readOnly attributes shall not be modified",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(100),
			EndLine:   l7643_7(100),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "readonly_not_modified",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
					"meta": map[string]any{
						"resourceType": "Invalid",
					},
				})
				r.RequireOK(err)

				if putResp.StatusCode == 200 && putResp.Body != nil {
					rt := GetString(putResp.Body, "meta", "resourceType")
					r.Check(rt == "User",
						fmt.Sprintf("after PUT with invalid meta.resourceType: got %q, want \"User\"", rt))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1766",
		Level:   ShallNot,
		Summary: "Immutable attributes shall not be updated after creation",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(105),
			EndLine:   l7643_7(107),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "immutable_not_updated",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
					"id":       "changed-id",
				})
				r.RequireOK(err)

				if putResp.StatusCode == 200 {
					r.Check(IDOf(putResp.Body) == id,
						fmt.Sprintf("id changed from %q to %q after PUT", id, IDOf(putResp.Body)))
				} else {
					r.Check(putResp.StatusCode == 400,
						FmtStatus("PUT with changed immutable id", putResp.StatusCode, 400))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-7-L1769",
		Level:   ShallNot,
		Summary: "writeOnly attribute values shall not be returned",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(108),
			EndLine:   l7643_7(111),
			EndCol:    25,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "write_only_not_returned",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				// password is writeOnly; it should not be returned.
				_, hasPassword := body["password"]
				r.Check(!hasPassword,
					"writeOnly attribute 'password' was returned in response")

				// Also check on GET.
				id := IDOf(body)
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					_, hasPassword = getResp.Body["password"]
					r.Check(!hasPassword,
						"GET returned writeOnly attribute 'password'")
				}
			},
		}},
	},
}
