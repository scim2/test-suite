package spec

var l7644_7 = rfcLines(rfc7644Sections, "26-section-7.txt")

var rfc7644_7 = []Requirement{
	// Section 7.2

	{
		ID:      "RFC7644-7.2-L4348",
		Level:   Must,
		Summary: "Clients and service providers must require transport-layer security",
		Source: Source{
			RFC:       7644,
			Section:   "7.2",
			StartLine: l7644_7(17),
			EndLine:   l7644_7(20),
			EndCol:    60,
		},
		Feature:  Core,
		Testable: false,
	},
	{
		ID:      "RFC7644-7.2-L4350",
		Level:   Must,
		Summary: "Service provider must support TLS 1.2",
		Source: Source{
			RFC:       7644,
			Section:   "7.2",
			StartLine: l7644_7(20),
			EndLine:   l7644_7(23),
			StartCol:  63,
		},
		Feature:  Core,
		Testable: false,
	},
}
