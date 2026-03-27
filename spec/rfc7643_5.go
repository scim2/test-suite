package spec

var l7643_5 = rfcLines(rfc7643Sections, "18-section-5.txt")

var rfc7643_5 = []Requirement{
	// Section 5

	{
		ID:      "RFC7643-5-L1581",
		Level:   Should,
		Summary: "authenticationSchemes should be publicly accessible without prior authentication",
		Source: Source{
			RFC:       7643,
			Section:   "5",
			StartLine: l7643_5(104),
			EndLine:   l7643_5(107),
			StartCol:  42,
			EndCol:    66,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "auth_schemes_present",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				if resp.StatusCode != 200 {
					r.Fatalf("GET /ServiceProviderConfig returned %d", resp.StatusCode)
				}

				authSchemes, _ := resp.Body["authenticationSchemes"].([]any)
				r.Check(len(authSchemes) > 0,
					"ServiceProviderConfig: authenticationSchemes is empty or missing")
			},
		}},
	},
}
