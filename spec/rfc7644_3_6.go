package spec

var l7644_3_6 = rfcLines(rfc7644Sections, "14-section-3.6.txt")

var rfc7644_3_6 = []Requirement{
	// Section 3.6

	{
		ID:      "RFC7644-3.6-L2642",
		Level:   Must,
		Summary: "Deleted resource must return 404 for all subsequent operations",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(3),
			EndLine:   l7644_3_6(7),
			StartCol:  50,
			EndCol:    38,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.6-L2644",
		Level:   Must,
		Summary: "Deleted resource must be omitted from future query results",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(6),
			EndLine:   l7644_3_6(7),
			StartCol:  34,
			EndCol:    38,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.6-L2657",
		Level:   Shall,
		Summary: "Successful DELETE returns 204 No Content",
		Source: Source{
			RFC:       7644,
			Section:   "3.6",
			StartLine: l7644_3_6(19),
			EndLine:   l7644_3_6(20),
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
	},
}
