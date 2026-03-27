package spec

var l7644_3_11 = rfcLines(rfc7644Sections, "19-section-3.11.txt")

var rfc7644_3_11 = []Requirement{
	// Section 3.11

	{
		ID:      "RFC7644-3.11-L3683",
		Level:   Must,
		Summary: "/Me response Location header must be the permanent resource location",
		Source: Source{
			RFC:       7644,
			Section:   "3.11",
			StartLine: l7644_3_11(16),
			EndLine:   l7644_3_11(19),
			StartCol:  66,
		},
		Feature:  Me,
		Testable: true,
	},
}
