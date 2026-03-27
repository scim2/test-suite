package spec

var l7644_5 = rfcLines(rfc7644Sections, "24-section-5.txt")

var rfc7644_5 = []Requirement{
	// Section 5

	{
		ID:      "RFC7644-5-L4214",
		Level:   Must,
		Summary: "Service providers must use PRECIS rules for userName and password comparison",
		Source: Source{
			RFC:       7644,
			Section:   "5",
			StartLine: l7644_5(7),
			EndLine:   l7644_5(13),
			StartCol:  16,
		},
		Feature:  Core,
		Testable: false,
	},
}
