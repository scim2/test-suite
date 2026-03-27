package spec

var l7644_3_9 = rfcLines(rfc7644Sections, "17-section-3.9.txt")

var rfc7644_3_9 = []Requirement{
	// Section 3.9

	{
		ID:      "RFC7644-3.9-L3596",
		Level:   Must,
		Summary: "attributes parameter overrides default list and must include minimum set",
		Source: Source{
			RFC:       7644,
			Section:   "3.9",
			StartLine: l7644_3_9(26),
			EndLine:   l7644_3_9(33),
		},
		Feature:  Core,
		Testable: true,
	},
}
