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
	},
}
