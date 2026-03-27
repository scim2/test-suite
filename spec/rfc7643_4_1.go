package spec

var l7643_4_1 = rfcLines(rfc7643Sections, "15-section-4.1.txt")

var rfc7643_4_1 = []Requirement{
	// Section 4.1.1

	{
		ID:      "RFC7643-4.1.1-L1035",
		Level:   Must,
		Summary: "Each User must include a non-empty userName value",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(15),
			EndLine:   l7643_4_1(18),
			StartCol:  51,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.1-L1186",
		Level:   ShallNot,
		Summary: "Cleartext or hashed password values shall not be returnable",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(165),
			EndLine:   l7643_4_1(167),
			StartCol:  59,
			EndCol:    22,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.1-L1222",
		Level:   MustNot,
		Summary: "Password value must not be returned in any form (returned is never)",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(201),
			EndLine:   l7643_4_1(204),
		},
		Feature:  Users,
		Testable: true,
	},

	// Section 4.1.2

	{
		ID:      "RFC7643-4.1.2-L1244",
		Level:   Should,
		Summary: "Email values should be specified per RFC 5321",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(223),
			EndLine:   l7643_4_1(227),
			EndCol:    36,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1281",
		Level:   Must,
		Summary: "Photo URI resource must be an image file, not a web page",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(260),
			EndLine:   l7643_4_1(263),
			EndCol:    15,
		},
		Feature:  Users,
		Testable: false,
	},
	{
		ID:      "RFC7643-4.1.2-L1322",
		Level:   Must,
		Summary: "Country value must be in ISO 3166-1 alpha-2 code format",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(300),
			EndLine:   l7643_4_1(303),
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1351",
		Level:   Must,
		Summary: "Group value sub-attribute must be the id of the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(321),
			EndLine:   l7643_4_1(336),
			StartCol:  63,
		},
		Feature:  Groups,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1354",
		Level:   Must,
		Summary: "Group membership changes must be applied via the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(333),
			EndLine:   l7643_4_1(336),
			StartCol:  22,
		},
		Feature:  Groups,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1378",
		Level:   Must,
		Summary: "X.509 certificate values must be base64 encoded DER",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(355),
			EndLine:   l7643_4_1(359),
			EndCol:    41,
		},
		Feature:  Users,
		Testable: true,
	},
	{
		ID:      "RFC7643-4.1.2-L1379",
		Level:   MustNot,
		Summary: "A single certificate value must not contain multiple certificates",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(359),
			EndLine:   l7643_4_1(361),
			StartCol:  44,
		},
		Feature:  Users,
		Testable: false,
	},
}
