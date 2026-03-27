package spec

var l7643_7 = rfcLines(rfc7643Sections, "20-section-7.txt")

var rfc7643_7 = []Requirement{
	// Section 7

	{
		ID:      "RFC7643-7-L1691",
		Level:   Must,
		Summary: "Service providers must specify the schema URI",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(31),
			EndLine:   l7643_7(33),
			EndCol:    51,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1700",
		Level:   Must,
		Summary: "Service providers must specify the schema name when applicable",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(40),
			EndLine:   l7643_7(41),
			StartCol:  42,
			EndCol:    63,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1748",
		Level:   Shall,
		Summary: "Server shall use case sensitivity when evaluating filters per caseExact",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(88),
			EndLine:   l7643_7(90),
			EndCol:    45,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1750",
		Level:   Shall,
		Summary: "Server shall preserve case for case-exact attribute values",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(90),
			EndLine:   l7643_7(92),
			StartCol:  48,
			EndCol:    19,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1759",
		Level:   ShallNot,
		Summary: "readOnly attributes shall not be modified",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(100),
			EndLine:   l7643_7(100),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1766",
		Level:   ShallNot,
		Summary: "Immutable attributes shall not be updated after creation",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(105),
			EndLine:   l7643_7(107),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-7-L1769",
		Level:   ShallNot,
		Summary: "writeOnly attribute values shall not be returned",
		Source: Source{
			RFC:       7643,
			Section:   "7",
			StartLine: l7643_7(108),
			EndLine:   l7643_7(111),
			EndCol:    25,
		},
		Feature:  Core,
		Testable: true,
	},
}
