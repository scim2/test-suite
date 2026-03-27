package spec

var l7643_1_1 = rfcLines(rfc7643Sections, "02-section-1.1.txt")

var rfc7643_1_1 = []Requirement{
	// Section 1.1

	{
		ID:      "RFC7643-1.1-L232",
		Level:   MustNot,
		Summary: "Quoted values in the spec must not include the quotes in protocol messages",
		Source: Source{
			RFC:       7643,
			Section:   "1.1",
			StartLine: l7643_1_1(25),
			EndLine:   l7643_1_1(27),
		},
		Feature:  Core,
		Testable: false,
	},
}
