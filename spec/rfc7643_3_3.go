package spec

var l7643_3_3 = rfcLines(rfc7643Sections, "13-section-3.3.txt")

var rfc7643_3_3 = []Requirement{
	// Section 3.3

	{
		ID:      "RFC7643-3.3-L976",
		Level:   Must,
		Summary: "The schemas attribute must contain at least one value (the base schema)",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: l7643_3_3(8),
			EndLine:   l7643_3_3(10),
			StartCol:  50,
			EndCol:    16,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "schemas_at_least_one",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				schemas := GetStringSlice(body, "schemas")
				r.Check(len(schemas) >= 1,
					"POST /Users: schemas must contain at least one value")
			},
		}},
	},
	{
		ID:      "RFC7643-3.3-L981",
		Level:   Shall,
		Summary: "Extension schema URI shall be used as JSON container for extension attributes",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: l7643_3_3(13),
			EndLine:   l7643_3_3(16),
			StartCol:  49,
			EndCol:    41,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "extension_container",
			Fn: func(r *Run) {
				body, resp := r.CreateTestResource(map[string]any{
					"schemas": []string{TestResourceSchema, TestExtSchema},
					TestExtSchema: map[string]any{
						"extString":   "ext-value",
						"extRequired": "required-value",
					},
				})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				extData, ok := body[TestExtSchema].(map[string]any)
				r.Check(ok && extData != nil,
					"extension attributes not found under schema URI key")
			},
		}},
	},
}
