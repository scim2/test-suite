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
	},
}
