package spec

var l7643_3_3 = rfcLines(rfc7643Sections, "13-section-3.3.txt")

var rfc7643_3_3 = []Requirement{
	// Section 3.3

	{
		ID:      "RFC7643-3.3-L976",
		Level:   Must,
		Summary: "The schemas attribute must contain at least one value (the base schema)",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: l7643_3_3(8),
			EndLine:   l7643_3_3(10),
			StartCol:  50,
			EndCol:    16,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3.3-L981",
		Level:   Shall,
		Summary: "Extension schema URI shall be used as JSON container for extension attributes",
		Source: Source{
			RFC:       7643,
			Section:   "3.3",
			StartLine: l7643_3_3(13),
			EndLine:   l7643_3_3(16),
			StartCol:  49,
			EndCol:    41,
		},
		Feature:  Core,
		Testable: true,
	},
}
