package spec

var l7644_3_13 = rfcLines(rfc7644Sections, "21-section-3.13.txt")

var rfc7644_3_13 = []Requirement{
	// Section 3.13

	{
		ID:      "RFC7644-3.13-L3956",
		Level:   Must,
		Summary: "When version is specified, service provider must use it or reject the request",
		Source: Source{
			RFC:       7644,
			Section:   "3.13",
			StartLine: l7644_3_13(7),
			EndLine:   l7644_3_13(12),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "if_match_wrong_version",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				putResp := r.DoWithHeaders(
					"PUT", "/Users/"+id,
					map[string]any{
						"schemas":  []string{UserSchema},
						"userName": userName,
					},
					map[string]string{"If-Match": "W/\"wrong-version\""},
				)

				r.Check(putResp.StatusCode == 412,
					FmtStatus("PUT with wrong If-Match", putResp.StatusCode, 412))
			},
		}},
	},
}
