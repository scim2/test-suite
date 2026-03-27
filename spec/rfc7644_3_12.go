package spec

import "fmt"

var l7644_3_12 = rfcLines(rfc7644Sections, "20-section-3.12.txt")

var rfc7644_3_12 = []Requirement{
	// Section 3.12

	{
		ID:      "RFC7644-3.12-L3707",
		Level:   Must,
		Summary: "Errors must be returned in JSON format in the response body",
		Source: Source{
			RFC:       7644,
			Section:   "3.12",
			StartLine: l7644_3_12(5),
			EndLine:   l7644_3_12(10),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "error_format",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Users/nonexistent-" + RandomSuffix())
				r.RequireOK(err)

				if resp.StatusCode == 404 && resp.Body != nil {
					r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:Error"),
						"404 error response missing Error schema")

					status, _ := resp.Body["status"].(string)
					r.Check(status == "404",
						fmt.Sprintf("error response status = %q, want \"404\"", status))
				}
			},
		}},
	},
}
