package spec

var l7644_3_3 = rfcLines(rfc7644Sections, "11-section-3.3.txt")

var rfc7644_3_3 = []Requirement{
	// Section 3.3

	{
		ID:      "RFC7644-3.3-L584",
		Level:   Shall,
		Summary: "readOnly attributes in POST request body shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(10),
			EndLine:   l7644_3_3(11),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L602",
		Level:   Shall,
		Summary: "Successful resource creation returns HTTP 201 Created",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(28),
			EndLine:   l7644_3_3(29),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L603",
		Level:   Should,
		Summary: "Response body should contain the newly created resource representation",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(30),
			EndLine:   l7644_3_3(31),
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L605",
		Level:   Shall,
		Summary: "Created resource URI shall be in Location header and meta.location",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(31),
			EndLine:   l7644_3_3(37),
			StartCol:  51,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.3-L625",
		Level:   Must,
		Summary: "Duplicate resource creation must return 409 Conflict with uniqueness scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.3",
			StartLine: l7644_3_3(50),
			EndLine:   l7644_3_3(54),
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.3.1

	{
		ID:      "RFC7644-3.3.1-L712",
		Level:   Shall,
		Summary: "meta.resourceType shall be set by service provider to match endpoint",
		Source: Source{
			RFC:       7644,
			Section:   "3.3.1",
			StartLine: l7644_3_3(138),
			EndLine:   l7644_3_3(142),
		},
		Feature:  Core,
		Testable: true,
	},
}
