package spec

import "fmt"

var l7643_2_5 = rfcLines(rfc7643Sections, "09-section-2.5.txt")

var rfc7643_2_5 = []Requirement{
	// Section 2.5

	{
		ID:      "RFC7643-2.5-L682",
		Level:   Shall,
		Summary: "Unassigned attributes, null, and empty array are equivalent in state",
		Source: Source{
			RFC:       7643,
			Section:   "2.5",
			StartLine: l7643_2_5(2),
			EndLine:   l7643_2_5(6),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "unassigned_equivalence",
			Fn: func(r *Run) {
				// Create user without displayName.
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				// displayName was not sent, so it should be absent or null.
				dn, hasDN := body["displayName"]
				r.Check(!hasDN || dn == nil || dn == "",
					fmt.Sprintf("unassigned displayName returned as %v", dn))
			},
		}},
	},
}
