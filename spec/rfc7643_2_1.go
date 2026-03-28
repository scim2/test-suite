package spec

var l7643_2_1 = rfcLines(rfc7643Sections, "05-section-2.1.txt")

var rfc7643_2_1 = []Requirement{
	// Section 2.1

	{
		ID:      "RFC7643-2.1-L366",
		Level:   Must,
		Summary: "Resources must be JSON and specify schema via the schemas attribute",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: l7643_2_1(11),
			EndLine:   l7643_2_1(13),
			StartCol:  26,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "schemas_present",
			Fn: func(r *Run) {
				for _, rt := range r.DiscoverResourceTypes() {
					r.Subtest(rt.Name, func(r *Run) {
						body, resp := r.CreateFuzzedResource(rt)
						if resp.StatusCode != 201 {
							r.Fatalf("POST %s returned %d", rt.Endpoint, resp.StatusCode)
						}

						r.Check(body != nil,
							"response is not a JSON object")

						schemas := GetStringSlice(body, "schemas")
						r.Check(len(schemas) > 0,
							"response missing schemas attribute")
					})
				}
			},
		}},
	},
	{
		ID:      "RFC7643-2.1-L369",
		Level:   Must,
		Summary: "Attribute names must conform to ATTRNAME ABNF rules",
		Source: Source{
			RFC:       7643,
			Section:   "2.1",
			StartLine: l7643_2_1(15),
			EndLine:   l7643_2_1(18),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "attribute_name_format",
			Fn: func(r *Run) {
				for _, rt := range r.DiscoverResourceTypes() {
					r.Subtest(rt.Name, func(r *Run) {
						body, resp := r.CreateFuzzedResource(rt)
						if resp.StatusCode != 201 {
							r.Fatalf("POST %s returned %d", rt.Endpoint, resp.StatusCode)
						}

						valid := true
						for key := range body {
							if IsSchemaURI(key) {
								continue
							}
							if !IsValidAttrName(key) {
								valid = false
								r.Logf("invalid attribute name: %q", key)
							}
						}
						r.Check(valid,
							"response contains attribute names not conforming to ATTRNAME ABNF")
					})
				}
			},
		}},
	},
}
