package spec

var RFC7642 = []Requirement{
	{
		ID:      "RFC7642-4-L944",
		Level:   Must,
		Summary: "Transport layer must guarantee data confidentiality (TLS)",
		Source: Source{
			RFC:       7642,
			Section:   "4",
			StartLine: 944,
			EndLine:   945,
		},
		Feature:  Core,
		Testable: false,
	},
}
