package spec

var l7644_2 = rfcLines(rfc7644Sections, "05-section-2.txt")

var rfc7644_2 = []Requirement{
	// Section 2

	{
		ID:      "RFC7644-2-L307",
		Level:   Shall,
		Summary: "Service provider shall indicate supported HTTP auth schemes via WWW-Authenticate",
		Source: Source{
			RFC:       7644,
			Section:   "2",
			StartLine: l7644_2(77),
			EndLine:   l7644_2(79),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "www_authenticate_header",
			Fn: func(r *Run) {
				rawClient := r.RawClient()
				resp, err := rawClient.Get("/Users")
				r.RequireOK(err)

				if resp.StatusCode == 401 {
					wwwAuth := resp.Header.Get("WWW-Authenticate")
					r.Check(wwwAuth != "",
						"401 response missing WWW-Authenticate header")
				} else {
					r.Check(true, "")
				}
			},
		}},
	},
	{
		ID:      "RFC7644-2-L311",
		Level:   Must,
		Summary: "Service provider must map authenticated client to access control policy",
		Source: Source{
			RFC:       7644,
			Section:   "2",
			StartLine: l7644_2(81),
			EndLine:   l7644_2(88),
		},
		Feature:  Core,
		Testable: false,
	},
}
