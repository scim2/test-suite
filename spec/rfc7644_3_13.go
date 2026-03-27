package spec

var l7644_3_13 = rfcLines(rfc7644Sections, "21-section-3.13.txt")

var rfc7644_3_13 = []Requirement{
	// Section 3.13

	{
		ID:      "RFC7644-3.13-L3956",
		Level:   Must,
		Summary: "When version is specified, service provider must use it or reject the request",
		Source: Source{
			RFC:       7644,
			Section:   "3.13",
			StartLine: l7644_3_13(7),
			EndLine:   l7644_3_13(12),
		},
		Feature:  Core,
		Testable: true,
	},
}
