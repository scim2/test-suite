package spec

var l7644_3_2 = rfcLines(rfc7644Sections, "10-section-3.2.txt")

var rfc7644_3_2 = []Requirement{
	// Section 3.2

	{
		ID:      "RFC7644-3.2-L490",
		Level:   MustNot,
		Summary: "PUT must not be used to create new resources",
		Source: Source{
			RFC:       7644,
			Section:   "3.2",
			StartLine: l7644_3_2(23),
			EndLine:   l7644_3_2(24),
			StartCol:  63,
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
}
