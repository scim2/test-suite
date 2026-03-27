package spec

var l7644_4 = rfcLines(rfc7644Sections, "23-section-4.txt")

var rfc7644_4 = []Requirement{
	// Section 4

	{
		ID:      "RFC7644-4-L4068",
		Level:   Shall,
		Summary: "/ServiceProviderConfig shall return ServiceProviderConfig schema",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: l7644_4(7),
			EndLine:   l7644_4(11),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "service_provider_config",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /ServiceProviderConfig", resp.StatusCode, 200))

				r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"),
					"ServiceProviderConfig missing expected schemas value")
			},
		}},
	},
	{
		ID:      "RFC7644-4-L4098",
		Level:   Shall,
		Summary: "/Schemas shall return all supported schemas in ListResponse format",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: l7644_4(37),
			EndLine:   l7644_4(42),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "schemas_endpoint",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Schemas")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /Schemas", resp.StatusCode, 200))

				r.Check(HasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
					"/Schemas response missing ListResponse schema")
			},
		}},
	},
	{
		ID:      "RFC7644-4-L4125",
		Level:   Shall,
		Summary: "Query parameters (filter, sort, pagination) shall be ignored on discovery endpoints",
		Source: Source{
			RFC:       7644,
			Section:   "4",
			StartLine: l7644_4(64),
			EndLine:   l7644_4(69),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "discovery_ignores_query_params",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ServiceProviderConfig?filter=invalid&sortBy=id&startIndex=999")
				r.RequireOK(err)

				r.Check(resp.StatusCode == 200,
					FmtStatus("GET /ServiceProviderConfig with query params", resp.StatusCode, 200))
			},
		}},
	},
}
