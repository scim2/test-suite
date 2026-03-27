package spec

var l7643_6 = rfcLines(rfc7643Sections, "19-section-6.txt")

var rfc7643_6 = []Requirement{
	// Section 6

	{
		ID:      "RFC7643-6-L1619",
		Level:   Must,
		Summary: "Service providers must specify the resource type name",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(18),
			EndLine:   l7643_6(19),
			StartCol:  32,
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1641",
		Level:   Must,
		Summary: "ResourceType schema URI must equal the id of the associated Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(38),
			EndLine:   l7643_6(41),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1650",
		Level:   Must,
		Summary: "Extension schema URI must equal the id of a Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(47),
			EndLine:   l7643_6(50),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-6-L1655",
		Level:   Must,
		Summary: "If schema extension is required, resource must include it and its required attributes",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(51),
			EndLine:   l7643_6(57),
		},
		Feature:  Core,
		Testable: true,
	},
}
