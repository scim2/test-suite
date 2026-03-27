package spec

var l7644_3_14 = rfcLines(rfc7644Sections, "22-section-3.14.txt")

var rfc7644_3_14 = []Requirement{
	// Section 3.14

	{
		ID:      "RFC7644-3.14-L3967",
		Level:   Must,
		Summary: "When supported, ETags must be specified as an HTTP header",
		Source: Source{
			RFC:       7644,
			Section:   "3.14",
			StartLine: l7644_3_14(5),
			EndLine:   l7644_3_14(9),
		},
		Feature:  ETag,
		Testable: true,
	},
}
