package spec

var l7644_3_8 = rfcLines(rfc7644Sections, "16-section-3.8.txt")

var rfc7644_3_8 = []Requirement{
	// Section 3.8

	{
		ID:      "RFC7644-3.8-L3537",
		Level:   Must,
		Summary: "Servers must accept and return JSON using UTF-8 encoding",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: l7644_3_8(3),
			EndLine:   l7644_3_8(5),
			EndCol:    19,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.8-L3551",
		Level:   Must,
		Summary: "Service providers must support Accept: application/scim+json",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: l7644_3_8(17),
			EndLine:   l7644_3_8(21),
		},
		Feature:  Core,
		Testable: true,
	},
}
