package spec

var l7643_2_4 = rfcLines(rfc7643Sections, "08-section-2.4.txt")

var rfc7643_2_4 = []Requirement{
	// Section 2.4

	{
		ID:      "RFC7643-2.4-L608",
		Level:   Shall,
		Summary: "Multi-valued attribute objects with sub-attributes shall be considered complex",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(9),
			EndLine:   l7643_2_4(11),
			EndCol:    56,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L612",
		Level:   Must,
		Summary: "Predefined multi-valued sub-attributes must be used with the meanings defined in spec",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(11),
			EndLine:   l7643_2_4(15),
			StartCol:  59,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7643-2.4-L634",
		Level:   Must,
		Summary: "The primary attribute value true must appear no more than once",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(36),
			EndLine:   l7643_2_4(38),
			StartCol:  35,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.4-L655",
		Level:   Should,
		Summary: "Service providers should canonicalize multi-valued attribute values when returning",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(58),
			EndLine:   l7643_2_4(61),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.4-L663",
		Level:   ShouldNot,
		Summary: "Should not return the same (type, value) combination more than once per attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.4",
			StartLine: l7643_2_4(63),
			EndLine:   l7643_2_4(67),
		},
		Feature:  Core,
		Testable: true,
	},
}
