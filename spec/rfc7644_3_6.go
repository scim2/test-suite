package spec

import (
	"fmt"
	"net/url"
)

var l7644_3_6 = rfcLines(rfc7644Sections, "14-section-3.6.txt")

var rfc7644_3_6 = []Requirement{
	// Section 3.6

	{
		ID:      "RFC7644-3.6-L2642",
		Level:   Must,
		Summary: "Deleted resource must return 404 for all subsequent operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(3),
			EndLine:   l7644_3_6(7),
			StartCol:  50,
			EndCol:    38,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "get_after_delete_returns_404",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-del-" + RandomSuffix(),
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d, need 201", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				if id == "" {
					r.Fatalf("setup: POST /Users returned no id")
				}

				delResp, err := r.Client.Delete("/Users/" + id)
				r.RequireOK(err)
				if delResp.StatusCode != 204 {
					r.Fatalf("setup: DELETE returned %d", delResp.StatusCode)
				}

				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				r.Check(getResp.StatusCode == 404,
					FmtStatus("GET /Users/"+id+" after DELETE", getResp.StatusCode, 404))
			},
		}},
	},
	{
		ID:      "RFC7644-3.6-L2644",
		Level:   Must,
		Summary: "Deleted resource must be omitted from future query results",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(6),
			EndLine:   l7644_3_6(7),
			StartCol:  34,
			EndCol:    38,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "deleted_omitted_from_query",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-delquery-" + RandomSuffix(),
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				userName, _ := resp.Body["userName"].(string)

				delResp, err := r.Client.Delete("/Users/" + id)
				r.RequireOK(err)
				if delResp.StatusCode != 204 {
					r.Fatalf("setup: DELETE returned %d", delResp.StatusCode)
				}

				filter := fmt.Sprintf("userName eq %q", userName)
				queryResp, err := r.Client.Get("/Users?filter=" + url.QueryEscape(filter))
				r.RequireOK(err)

				tr, _ := queryResp.Body["totalResults"].(float64)
				r.Check(tr == 0,
					fmt.Sprintf("query for deleted user returned totalResults=%v, want 0", queryResp.Body["totalResults"]))
			},
		}},
	},
	{
		ID:      "RFC7644-3.6-L2657",
		Level:   Shall,
		Summary: "Successful DELETE returns 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(19),
			EndLine:   l7644_3_6(20),
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "delete_returns_204",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-del-" + RandomSuffix(),
				})
				r.RequireOK(err)
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d, need 201", resp.StatusCode)
				}
				id := IDOf(resp.Body)
				if id == "" {
					r.Fatalf("setup: POST /Users returned no id")
				}

				delResp, err := r.Client.Delete("/Users/" + id)
				r.RequireOK(err)

				r.Check(delResp.StatusCode == 204,
					FmtStatus("DELETE /Users/"+id, delResp.StatusCode, 204))
			},
		}},
	},
}
