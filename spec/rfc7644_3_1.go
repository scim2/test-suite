package spec

var l7644_3_1 = rfcLines(rfc7644Sections, "09-section-3.1.txt")

var rfc7644_3_1 = []Requirement{
	// Section 3.1

	{
		ID:      "RFC7644-3.1-L436",
		Level:   Must,
		Summary: "SCIM message schemas URI prefix must begin with urn:ietf:params:scim:api:",
		Source: Source{
			RFC:       7644,
			Section:   "3.1",
			StartLine: l7644_3_1(34),
			EndLine:   l7644_3_1(37),
			StartCol:  63,
			EndCol:    42,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7644-3.1-L439",
		Level:   ShallNot,
		Summary: "SCIM message schemas shall not be discoverable via /Schemas",
		Source: Source{
			RFC:       7644,
			Section:   "3.1",
			StartLine: l7644_3_1(38),
			EndLine:   l7644_3_1(40),
		},
		Feature:  Core,
		Testable: true,
	},
}
