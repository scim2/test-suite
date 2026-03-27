package spec

var l7643_2_1 = rfcLines(rfc7643Sections, "05-section-2.1.txt")

var rfc7643_2_1 = []Requirement{
	// Section 2.1

	{
		ID:      "RFC7643-2.1-L366",
		Level:   Must,
		Summary: "Resources must be JSON and specify schema via the schemas attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: l7643_2_1(11),
			EndLine:   l7643_2_1(13),
			StartCol:  26,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.1-L369",
		Level:   Must,
		Summary: "Attribute names must conform to ATTRNAME ABNF rules",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: l7643_2_1(15),
			EndLine:   l7643_2_1(18),
		},
		Feature:  Core,
		Testable: true,
	},
}
