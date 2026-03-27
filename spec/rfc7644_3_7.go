package spec

var l7644_3_7 = rfcLines(rfc7644Sections, "15-section-3.7.txt")

var rfc7644_3_7 = []Requirement{
	// Section 3.7

	{
		ID:      "RFC7644-3.7-L2775",
		Level:   Must,
		Summary: "Successful bulk job must return 200 OK",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(81),
			EndLine:   l7644_3_7(83),
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2779",
		Level:   Must,
		Summary: "Service provider must continue performing changes and disregard partial failures",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(85),
			EndLine:   l7644_3_7(86),
			EndCol:    43,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2789",
		Level:   Must,
		Summary: "Service provider must return the same bulkId with newly created resource",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(95),
			EndLine:   l7644_3_7(98),
			StartCol:  36,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7-L2796",
		Level:   Must,
		Summary: "Optimized bulk processing must preserve client intent and achieve same result",
		Source: Source{
			RFC:       7644,
			Section:   "3.7",
			StartLine: l7644_3_7(101),
			EndLine:   l7644_3_7(113),
			StartCol:  68,
			EndCol:    31,
		},
		Feature:  Bulk,
		Testable: false,
	},

	// Section 3.7.1

	{
		ID:      "RFC7644-3.7.1-L2814",
		Level:   Must,
		Summary: "Service provider must try to resolve circular cross-references in bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.1",
			StartLine: l7644_3_7(120),
			EndLine:   l7644_3_7(122),
			EndCol:    62,
		},
		Feature:  Bulk,
		Testable: true,
	},

	// Section 3.7.3

	{
		ID:      "RFC7644-3.7.3-L3201",
		Level:   Must,
		Summary: "Bulk response must include result of all processed operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: l7644_3_7(507),
			EndLine:   l7644_3_7(512),
			EndCol:    33,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7.3-L3206",
		Level:   Must,
		Summary: "Bulk status must include HTTP response code",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.3",
			StartLine: l7644_3_7(511),
			EndLine:   l7644_3_7(516),
		},
		Feature:  Bulk,
		Testable: true,
	},

	// Section 3.7.4

	{
		ID:      "RFC7644-3.7.4-L3495",
		Level:   Must,
		Summary: "Service provider must define max operations and max payload size for bulk",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: l7644_3_7(801),
			EndLine:   l7644_3_7(804),
			EndCol:    50,
		},
		Feature:  Bulk,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.7.4-L3499",
		Level:   Must,
		Summary: "Exceeding bulk limits must return 413 Payload Too Large",
		Source: Source{
			RFC:       7644,
			Section:   "3.7.4",
			StartLine: l7644_3_7(804),
			EndLine:   l7644_3_7(807),
			StartCol:  49,
		},
		Feature:  Bulk,
		Testable: true,
	},
}
