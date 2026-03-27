package spec

var l7643_3 = rfcLines(rfc7643Sections, "10-section-3.txt")

var rfc7643_3 = []Requirement{
	// Section 3

	{
		ID:      "RFC7643-3-L709",
		Level:   Must,
		Summary: "All schema representations must include a non-empty schemas array",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: l7643_3(21),
			EndLine:   l7643_3(23),
			StartCol:  20,
			EndCol:    21,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L711",
		Level:   Must,
		Summary: "Schemas attribute must only contain values defined as schema and schemaExtensions",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: l7643_3(23),
			EndLine:   l7643_3(25),
			StartCol:  24,
			EndCol:    40,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L713",
		Level:   MustNot,
		Summary: "Duplicate values must not be included in the schemas attribute",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: l7643_3(25),
			EndLine:   l7643_3(26),
			StartCol:  43,
			EndCol:    15,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-3-L714",
		Level:   MustNot,
		Summary: "Value order in schemas attribute must not impact behavior",
		Source: Source{
			RFC:       7643,
			Section:   "3",
			StartLine: l7643_3(26),
			EndLine:   l7643_3(27),
			StartCol:  18,
		},
		Feature:  Core,
		Testable: false,
	},
}
