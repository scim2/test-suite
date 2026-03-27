package spec

import "strings"

var l7644_3_14 = rfcLines(rfc7644Sections, "22-section-3.14.txt")

var rfc7644_3_14 = []Requirement{
	// Section 3.14

	{
		ID:      "RFC7644-3.14-L3967",
		Level:   Must,
		Summary: "When supported, ETags must be specified as an HTTP header",
		Source: Source{
			RFC:       7644,
			Section:   "3.14",
			StartLine: l7644_3_14(5),
			EndLine:   l7644_3_14(9),
		},
		Feature:  ETag,
		Testable: true,
		Tests: []Test{{
			Name: "etag_in_http_header",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				id := IDOf(body)
				getResp, err := r.Client.Get("/Users/" + id)
				r.RequireOK(err)

				etag := getResp.Header.Get("ETag")
				if etag != "" {
					r.Check(strings.HasPrefix(etag, "W/") || strings.HasPrefix(etag, "\""),
						"ETag header present but not a valid entity-tag")
				}
			},
		}},
	},
}
