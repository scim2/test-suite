package spec

var l7644_3_4 = rfcLines(rfc7644Sections, "12-section-3.4.txt")

var rfc7644_3_4 = []Requirement{
	// Section 3.4.2

	{
		ID:      "RFC7644-3.4.2-L817",
		Level:   Must,
		Summary: "Query responses must use ListResponse schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: l7644_3_4(101),
			EndLine:   l7644_3_4(102),
			EndCol:    55,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2-L847",
		Level:   Shall,
		Summary: "Empty query result returns 200 with totalResults 0",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2",
			StartLine: l7644_3_4(131),
			EndLine:   l7644_3_4(132),
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.2.1

	{
		ID:      "RFC7644-3.4.2.1-L912",
		Level:   Shall,
		Summary: "Too many results shall return 400 with scimType tooMany",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: l7644_3_4(194),
			EndLine:   l7644_3_4(198),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.1-L918",
		Level:   Must,
		Summary: "Multi-resource-type queries must process filters per Section 3.4.2.2",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.1",
			StartLine: l7644_3_4(200),
			EndLine:   l7644_3_4(202),
		},
		Feature:  Filter,
		Testable: true,
	},

	// Section 3.4.2.2

	{
		ID:      "RFC7644-3.4.2.2-L944",
		Level:   Must,
		Summary: "Filter parameter must contain at least one valid expression",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(228),
			EndLine:   l7644_3_4(231),
			EndCol:    53,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1001",
		Level:   Shall,
		Summary: "Boolean/Binary gt comparison shall fail with 400 invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(285),
			EndLine:   l7644_3_4(288),
			EndCol:    46,
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1127",
		Level:   Must,
		Summary: "SCIM filters must conform to the ABNF rules",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(411),
			EndLine:   l7644_3_4(412),
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1197",
		Level:   Must,
		Summary: "Complex attributes must use fully qualified sub-attribute in filters",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(481),
			EndLine:   l7644_3_4(483),
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1208",
		Level:   Must,
		Summary: "Unrecognized filter operations must return 400 with invalidFilter",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(492),
			EndLine:   l7644_3_4(498),
		},
		Feature:  Filter,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.2-L1217",
		Level:   Shall,
		Summary: "String filter comparisons shall respect caseExact characteristic",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.2",
			StartLine: l7644_3_4(500),
			EndLine:   l7644_3_4(502),
		},
		Feature:  Filter,
		Testable: true,
	},

	// Section 3.4.2.3

	{
		ID:      "RFC7644-3.4.2.3-L1319",
		Level:   Shall,
		Summary: "sortOrder shall default to ascending when sortBy is provided",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: l7644_3_4(600),
			EndLine:   l7644_3_4(603),
			EndCol:    33,
		},
		Feature:  Sort,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.3-L1321",
		Level:   Must,
		Summary: "sortOrder must sort according to attribute type (case-insensitive or case-exact)",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.3",
			StartLine: l7644_3_4(603),
			EndLine:   l7644_3_4(610),
			StartCol:  36,
		},
		Feature:  Sort,
		Testable: true,
	},

	// Section 3.4.2.4

	{
		ID:      "RFC7644-3.4.2.4-L1358",
		Level:   Shall,
		Summary: "startIndex value less than 1 shall be interpreted as 1",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: l7644_3_4(641),
			EndLine:   l7644_3_4(643),
			StartCol:  40,
			EndCol:    35,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.4-L1363",
		Level:   MustNot,
		Summary: "Service provider must not return more results than count specifies",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.4",
			StartLine: l7644_3_4(647),
			EndLine:   l7644_3_4(648),
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.2.5

	{
		ID:      "RFC7644-3.4.2.5-L1438",
		Level:   Must,
		Summary: "attributes and excludedAttributes parameters must be supported",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: l7644_3_4(721),
			EndLine:   l7644_3_4(723),
			StartCol:  31,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.4.2.5-L1449",
		Level:   Shall,
		Summary: "excludedAttributes shall have no effect on returned:always attributes",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.2.5",
			StartLine: l7644_3_4(731),
			EndLine:   l7644_3_4(736),
			EndCol:    54,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.4.3

	{
		ID:      "RFC7644-3.4.3-L1476",
		Level:   Must,
		Summary: "POST search requests must use SearchRequest schema URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.4.3",
			StartLine: l7644_3_4(760),
			EndLine:   l7644_3_4(762),
		},
		Feature:  Core,
		Testable: true,
	},
}
