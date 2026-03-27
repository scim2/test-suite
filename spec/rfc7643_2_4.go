package spec

import "fmt"

var l7643_2_4 = rfcLines(rfc7643Sections, "08-section-2.4.txt")

var rfc7643_2_4 = []Requirement{
	// Section 2.4

	{
		ID:      "RFC7643-2.4-L608",
		Level:   Shall,
		Summary: "Multi-valued attribute objects with sub-attributes shall be considered complex",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(9),
			EndLine:   l7643_2_4(11),
			EndCol:    56,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L612",
		Level:   Must,
		Summary: "Predefined multi-valued sub-attributes must be used with the meanings defined in spec",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(11),
			EndLine:   l7643_2_4(15),
			StartCol:  59,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L634",
		Level:   Must,
		Summary: "The primary attribute value true must appear no more than once",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(36),
			EndLine:   l7643_2_4(38),
			StartCol:  35,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "primary_uniqueness",
			Fn: func(r *Run) {
				// Create user with two emails, both marked primary.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-primary-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "a@example.com", "type": "work", "primary": true},
						{"value": "b@example.com", "type": "home", "primary": true},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode == 201 && resp.Body != nil {
					emails, _ := resp.Body["emails"].([]any)
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
						fmt.Sprintf("emails has %d primary=true values, want at most 1", primaryCount))
				} else if resp.StatusCode == 400 {
					// Server correctly rejected the input.
					r.Check(true, "")
				} else {
					r.Check(false,
						FmtStatus("POST /Users with duplicate primary", resp.StatusCode, 400))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-2.4-L655",
		Level:   Should,
		Summary: "Service providers should canonicalize multi-valued attribute values when returning",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(58),
			EndLine:   l7643_2_4(61),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "multivalued_canonicalize",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-canon-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "TEST@EXAMPLE.COM", "type": "work"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode == 201 && resp.Body != nil {
					emails, _ := resp.Body["emails"].([]any)
					if len(emails) > 0 {
						first, _ := emails[0].(map[string]any)
						v, _ := first["value"].(string)
						r.Check(v != "",
							"email value was empty after canonicalization")
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7643-2.4-L663",
		Level:   ShouldNot,
		Summary: "Should not return the same (type, value) combination more than once per attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(63),
			EndLine:   l7643_2_4(67),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "no_duplicate_type_value",
			Fn: func(r *Run) {
				// Create user with duplicate (type, value) email entries.
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-dedup-" + RandomSuffix(),
					"emails": []map[string]any{
						{"value": "dup@example.com", "type": "work"},
						{"value": "dup@example.com", "type": "work"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode == 201 && resp.Body != nil {
					emails, _ := resp.Body["emails"].([]any)

					type tv struct{ t, v string }
					seen := make(map[tv]bool)
					hasDupes := false
					for _, e := range emails {
						em, ok := e.(map[string]any)
						if !ok {
							continue
						}
						key := tv{
							t: fmt.Sprintf("%v", em["type"]),
							v: fmt.Sprintf("%v", em["value"]),
						}
						if seen[key] {
							hasDupes = true
						}
						seen[key] = true
					}
					r.Check(!hasDupes,
						fmt.Sprintf("emails contains duplicate (type, value) pairs: %v", emails))
				} else if resp.StatusCode == 400 {
					// Server correctly rejected duplicate input.
					r.Check(true, "")
				} else {
					r.Check(false,
						FmtStatus("POST /Users with duplicate (type,value)", resp.StatusCode, 400))
				}
			},
		}},
	},
}
