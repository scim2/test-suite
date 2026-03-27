package spec

var l7644_1_3 = rfcLines(rfc7644Sections, "04-section-1.3.txt")

var rfc7644_1_3 = []Requirement{
	// Section 1.3

	{
		ID:      "RFC7644-1.3-L201",
		Level:   MustNot,
		Summary: "Base URI must not contain a query string",
		Source: Source{
			RFC:       7644,
			Section:   "1.3",
			StartLine: l7644_1_3(7),
			EndLine:   l7644_1_3(9),
		},
		Feature:  Core,
		Testable: false,
	},
}
