package spec

import (
	"fmt"
	"net/url"
)

var l7644_3_9 = rfcLines(rfc7644Sections, "17-section-3.9.txt")

var rfc7644_3_9 = []Requirement{
	// Section 3.9

	{
		ID:      "RFC7644-3.9-L3596",
		Level:   Must,
		Summary: "attributes parameter overrides default list and must include minimum set",
		Source: Source{
			RFC:       7644,
			Section:   "3.9",
			StartLine: l7644_3_9(26),
			EndLine:   l7644_3_9(33),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{
			{
				Name: "select_on_list",
				Fn: func(r *Run) {
					r.CreateUser()

					resp, err := r.Client.Get("/Users?attributes=" + url.QueryEscape("userName"))
					r.RequireOK(err)

					if resp.StatusCode != 200 {
						r.Fatalf("GET /Users?attributes=userName returned %d", resp.StatusCode)
					}

					resources, _ := resp.Body["Resources"].([]any)
					if len(resources) == 0 {
						r.Skipf("no resources returned")
					}

					first, ok := resources[0].(map[string]any)
					if !ok {
						r.Fatalf("first resource is not a JSON object")
					}

					r.Check(IDOf(first) != "",
						"list with attributes=userName: resource missing id")

					userName, _ := first["userName"].(string)
					r.Check(userName != "",
						"list with attributes=userName: resource missing userName")

					_, hasDisplayName := first["displayName"]
					r.Check(!hasDisplayName,
						fmt.Sprintf("list with attributes=userName: displayName should not be returned, got %v", first["displayName"]))
				},
			},
			{
				Name: "select_username_includes_minimum_set",
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

					r.Check(IDOf(qResp.Body) != "",
						"GET with attributes=userName: id missing (should always be returned)")

					schemas := GetStringSlice(qResp.Body, "schemas")
					r.Check(len(schemas) > 0,
						"GET with attributes=userName: schemas missing (should always be returned)")
				},
			},
		},
	},
}
