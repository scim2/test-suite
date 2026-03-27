package spec

var l7644_3_5 = rfcLines(rfc7644Sections, "13-section-3.5.txt")

var rfc7644_3_5 = []Requirement{
	// Section 3.5

	{
		ID:      "RFC7644-3.5-L1607",
		Level:   Must,
		Summary: "Implementers must support HTTP PUT",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: l7644_3_5(3),
			EndLine:   l7644_3_5(5),
			EndCol:    31,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5-L1609",
		Level:   Should,
		Summary: "Implementers should support HTTP PATCH for partial modifications",
		Source: Source{
			RFC:       7644,
			Section:   "3.5",
			StartLine: l7644_3_5(5),
			EndLine:   l7644_3_5(9),
			StartCol:  34,
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.1

	{
		ID:      "RFC7644-3.5.1-L1637",
		Level:   MustNot,
		Summary: "HTTP PUT must not be used to create new resources",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(33),
			EndLine:   l7644_3_5(34),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1644",
		Level:   Shall,
		Summary: "readWrite/writeOnly values provided in PUT shall replace existing values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(41),
			EndLine:   l7644_3_5(42),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1660",
		Level:   Must,
		Summary: "Immutable attribute input must match existing values or return 400",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(56),
			EndLine:   l7644_3_5(60),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1665",
		Level:   Shall,
		Summary: "readOnly values provided in PUT shall be ignored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(62),
			EndLine:   l7644_3_5(62),
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.1-L1667",
		Level:   Must,
		Summary: "Required attributes must be specified in PUT request",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.1",
			StartLine: l7644_3_5(64),
			EndLine:   l7644_3_5(65),
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 3.5.2

	{
		ID:      "RFC7644-3.5.2-L1805",
		Level:   Must,
		Summary: "PATCH body must contain schemas attribute with PatchOp URI",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(202),
			EndLine:   l7644_3_5(203),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1808",
		Level:   Must,
		Summary: "PATCH body must contain Operations array with one or more operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(205),
			EndLine:   l7644_3_5(207),
			EndCol:    14,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1810",
		Level:   Must,
		Summary: "Each PATCH operation must have exactly one op member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(207),
			EndLine:   l7644_3_5(211),
			StartCol:  17,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1886",
		Level:   Must,
		Summary: "Each PATCH operation must be compatible with attribute mutability and schema",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(283),
			EndLine:   l7644_3_5(285),
			EndCol:    16,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1888",
		Level:   MustNot,
		Summary: "Client must not modify readOnly or immutable attributes via PATCH",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(285),
			EndLine:   l7644_3_5(288),
			StartCol:  19,
			EndCol:    21,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1927",
		Level:   Shall,
		Summary: "Setting primary to true shall auto-set primary to false for other values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(323),
			EndLine:   l7644_3_5(326),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1931",
		Level:   Shall,
		Summary: "PATCH request shall be treated as atomic",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(328),
			EndLine:   l7644_3_5(331),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1933",
		Level:   Must,
		Summary: "On PATCH operation error, original resource must be restored",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(329),
			EndLine:   l7644_3_5(331),
			StartCol:  24,
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2-L1939",
		Level:   Must,
		Summary: "Successful PATCH must return 200 with resource body or 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2",
			StartLine: l7644_3_5(336),
			EndLine:   l7644_3_5(341),
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.1

	{
		ID:      "RFC7644-3.5.2.1-L1972",
		Level:   Must,
		Summary: "PATCH add operation must contain a value member",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(369),
			EndLine:   l7644_3_5(372),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.1-L1988",
		Level:   Shall,
		Summary: "Add to complex attribute shall specify sub-attributes in value parameter",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(384),
			EndLine:   l7644_3_5(385),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.1-L2004",
		Level:   ShallNot,
		Summary: "Add of already-existing value shall not change the modify timestamp",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.1",
			StartLine: l7644_3_5(398),
			EndLine:   l7644_3_5(402),
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.2

	{
		ID:      "RFC7644-3.5.2.2-L2146",
		Level:   Shall,
		Summary: "Removing single-value attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(542),
			EndLine:   l7644_3_5(544),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.2-L2151",
		Level:   Shall,
		Summary: "Removing all values of multi-valued attribute shall make it unassigned",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(546),
			EndLine:   l7644_3_5(548),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.2-L2167",
		Level:   Shall,
		Summary: "Removing required or readOnly attribute shall return error with mutability scimType",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.2",
			StartLine: l7644_3_5(563),
			EndLine:   l7644_3_5(567),
		},
		Feature:  Patch,
		Testable: true,
	},

	// Section 3.5.2.3

	{
		ID:      "RFC7644-3.5.2.3-L2376",
		Level:   Shall,
		Summary: "Replace on non-existent attribute shall treat as add",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(772),
			EndLine:   l7644_3_5(773),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.3-L2387",
		Level:   Shall,
		Summary: "Replace with valuePath filter shall replace all matching record values",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(781),
			EndLine:   l7644_3_5(784),
		},
		Feature:  Patch,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.5.2.3-L2396",
		Level:   Shall,
		Summary: "Replace with unmatched valuePath filter shall return 400 noTarget",
		Source: Source{
			RFC:       7644,
			Section:   "3.5.2.3",
			StartLine: l7644_3_5(791),
			EndLine:   l7644_3_5(795),
		},
		Feature:  Patch,
		Testable: true,
	},
}
