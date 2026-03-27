package spec

var l7643_9 = rfcLines(rfc7643Sections, "22-section-9.txt")

var rfc7643_9 = []Requirement{
	// Section 9.2

	{
		ID:      "RFC7643-9.2-L5117",
		Level:   MustNot,
		Summary: "Password values must not be stored in cleartext form",
		Source: Source{
			RFC:       7643,
			Section:   "9.2",
			StartLine: l7643_9(12),
			EndLine:   l7643_9(15),
		},
		Feature:  Core,
		Testable: false,
	},
}
