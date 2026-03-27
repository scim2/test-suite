package spec

var l7643_2_3 = rfcLines(rfc7643Sections, "07-section-2.3.txt")

var rfc7643_2_3 = []Requirement{
	// Section 2.3

	{
		ID:      "RFC7643-2.3-L443",
		Level:   ShouldNot,
		Summary: "Extensions should not introduce new data types",
		Source: Source{
			RFC:       7643,
			Section:   "2.3",
			StartLine: l7643_2_3(6),
			EndLine:   l7643_2_3(7),
			StartCol:  24,
		},
		Feature:  Core,
		Testable: false,
	},

	// Section 2.3.4

	{
		ID:      "RFC7643-2.3.4-L521",
		Level:   MustNot,
		Summary: "Integer values must not contain fractional or exponent parts",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.4",
			StartLine: l7643_2_3(82),
			EndLine:   l7643_2_3(84),
			StartCol:  58,
			EndCol:    64,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.5

	{
		ID:      "RFC7643-2.3.5-L527",
		Level:   Must,
		Summary: "DateTime values must be valid xsd:dateTime with both date and time",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.5",
			StartLine: l7643_2_3(89),
			EndLine:   l7643_2_3(91),
			StartCol:  52,
			EndCol:    59,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.5-L531",
		Level:   Must,
		Summary: "DateTime JSON values must conform to XML constraints",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.5",
			StartLine: l7643_2_3(94),
			EndLine:   l7643_2_3(96),
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.6

	{
		ID:      "RFC7643-2.3.6-L537",
		Level:   Must,
		Summary: "Binary attribute values must be base64 encoded per RFC 4648",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.6",
			StartLine: l7643_2_3(100),
			EndLine:   l7643_2_3(101),
			StartCol:  28,
			EndCol:    39,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.7

	{
		ID:      "RFC7643-2.3.7-L552",
		Level:   Must,
		Summary: "Reference values must be absolute or relative URIs",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: l7643_2_3(115),
			EndLine:   l7643_2_3(116),
			EndCol:    12,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L555",
		Level:   Must,
		Summary: "Base URI for relative URI resolution must include all components up to endpoint URI",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: l7643_2_3(117),
			EndLine:   l7643_2_3(133),
			StartCol:  31,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L577",
		Level:   Must,
		Summary: "Reference URI must be HTTP-addressable",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: l7643_2_3(140),
			EndLine:   l7643_2_3(140),
			EndCol:    59,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.7-L578",
		Level:   Must,
		Summary: "GET on a reference URI must return the target resource or appropriate HTTP code",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.7",
			StartLine: l7643_2_3(140),
			EndLine:   l7643_2_3(142),
			StartCol:  62,
			EndCol:    56,
		},
		Feature:  Core,
		Testable: true,
	},

	// Section 2.3.8

	{
		ID:      "RFC7643-2.3.8-L592",
		Level:   MustNot,
		Summary: "Servers and clients must not require or expect attributes in any specific order",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.8",
			StartLine: l7643_2_3(154),
			EndLine:   l7643_2_3(157),
			StartCol:  29,
			EndCol:    25,
		},
		Feature:  Core,
		Testable: true,
	},
	{
		ID:      "RFC7643-2.3.8-L595",
		Level:   MustNot,
		Summary: "Complex attributes must not contain sub-attributes that are themselves complex",
		Source: Source{
			RFC:       7643,
			Section:   "2.3.8",
			StartLine: l7643_2_3(158),
			EndLine:   l7643_2_3(159),
			StartCol:  18,
		},
		Feature:  Core,
		Testable: true,
	},
}
