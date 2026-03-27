package spec

var l7643_2_5 = rfcLines(rfc7643Sections, "09-section-2.5.txt")

var rfc7643_2_5 = []Requirement{
	// Section 2.5

	{
		ID:      "RFC7643-2.5-L682",
		Level:   Shall,
		Summary: "Unassigned attributes, null, and empty array are equivalent in state",
		Source: Source{
			RFC:       7643,
			Section:   "2.5",
			StartLine: l7643_2_5(2),
			EndLine:   l7643_2_5(6),
		},
		Feature:  Core,
		Testable: true,
	},
}
