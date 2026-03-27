package spec

import "fmt"

var l7643_6 = rfcLines(rfc7643Sections, "19-section-6.txt")

var rfc7643_6 = []Requirement{
	// Section 6

	{
		ID:      "RFC7643-6-L1619",
		Level:   Must,
		Summary: "Service providers must specify the resource type name",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(18),
			EndLine:   l7643_6(19),
			StartCol:  32,
			EndCol:    48,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "resource_type_has_name",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ResourceTypes")
				r.RequireOK(err)

				if resp.StatusCode != 200 || resp.Body == nil {
					r.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
				}

				resources, _ := resp.Body["Resources"].([]any)
				for _, res := range resources {
					rt, ok := res.(map[string]any)
					if !ok {
						continue
					}
					name, _ := rt["name"].(string)
					r.Check(name != "",
						"ResourceType missing name")
				}
			},
		}},
	},
	{
		ID:      "RFC7643-6-L1641",
		Level:   Must,
		Summary: "ResourceType schema URI must equal the id of the associated Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(38),
			EndLine:   l7643_6(41),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{
			{
				Name: "resource_type_has_schema",
				Fn: func(r *Run) {
					resp, err := r.Client.Get("/ResourceTypes")
					r.RequireOK(err)

					if resp.StatusCode != 200 || resp.Body == nil {
						r.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
					}

					resources, _ := resp.Body["Resources"].([]any)
					for _, res := range resources {
						rt, ok := res.(map[string]any)
						if !ok {
							continue
						}
						name, _ := rt["name"].(string)
						schema, _ := rt["schema"].(string)
						r.Check(schema != "",
							fmt.Sprintf("ResourceType %q: missing schema URI", name))
					}
				},
			},
			{
				Name: "schema_uri_in_schemas_endpoint",
				Fn: func(r *Run) {
					resp, err := r.Client.Get("/ResourceTypes")
					r.RequireOK(err)

					if resp.StatusCode != 200 || resp.Body == nil {
						r.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
					}

					schemasResp, err := r.Client.Get("/Schemas")
					r.RequireOK(err)

					schemaIDs := make(map[string]bool)
					if schemasResp.StatusCode == 200 && schemasResp.Body != nil {
						schemaResources, _ := schemasResp.Body["Resources"].([]any)
						for _, sr := range schemaResources {
							s, ok := sr.(map[string]any)
							if !ok {
								continue
							}
							if id, ok := s["id"].(string); ok {
								schemaIDs[id] = true
							}
						}
					}

					resources, _ := resp.Body["Resources"].([]any)
					for _, res := range resources {
						rt, ok := res.(map[string]any)
						if !ok {
							continue
						}
						name, _ := rt["name"].(string)
						schemaURI, _ := rt["schema"].(string)

						if len(schemaIDs) > 0 && schemaURI != "" {
							r.Check(schemaIDs[schemaURI],
								fmt.Sprintf("ResourceType %q: schema %q not found in /Schemas", name, schemaURI))
						}
					}
				},
			},
		},
	},
	{
		ID:      "RFC7643-6-L1650",
		Level:   Must,
		Summary: "Extension schema URI must equal the id of a Schema resource",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(47),
			EndLine:   l7643_6(50),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "extension_schema_in_schemas_endpoint",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ResourceTypes")
				r.RequireOK(err)

				if resp.StatusCode != 200 || resp.Body == nil {
					r.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
				}

				schemasResp, err := r.Client.Get("/Schemas")
				r.RequireOK(err)

				schemaIDs := make(map[string]bool)
				if schemasResp.StatusCode == 200 && schemasResp.Body != nil {
					schemaResources, _ := schemasResp.Body["Resources"].([]any)
					for _, sr := range schemaResources {
						s, ok := sr.(map[string]any)
						if !ok {
							continue
						}
						if id, ok := s["id"].(string); ok {
							schemaIDs[id] = true
						}
					}
				}

				resources, _ := resp.Body["Resources"].([]any)
				for _, res := range resources {
					rt, ok := res.(map[string]any)
					if !ok {
						continue
					}
					name, _ := rt["name"].(string)

					exts, _ := rt["schemaExtensions"].([]any)
					for _, ext := range exts {
						e, ok := ext.(map[string]any)
						if !ok {
							continue
						}
						extSchema, _ := e["schema"].(string)
						if extSchema != "" && len(schemaIDs) > 0 {
							r.Check(schemaIDs[extSchema],
								fmt.Sprintf("ResourceType %q: extension schema %q not found in /Schemas", name, extSchema))
						}
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7643-6-L1655",
		Level:   Must,
		Summary: "If schema extension is required, resource must include it and its required attributes",
		Source: Source{
			RFC:       7643,
			Section:   "6",
			StartLine: l7643_6(51),
			EndLine:   l7643_6(57),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "required_extension",
			Fn: func(r *Run) {
				// Try to create a TestResource WITHOUT the required extension.
				resp, err := r.Client.Post("/TestResources", map[string]any{
					"schemas":    []string{TestResourceSchema},
					"identifier": "scim-test-noext-" + RandomSuffix(),
				})
				r.RequireOK(err)

				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/TestResources/" + id) })
					}
				}

				r.Check(resp.StatusCode == 400,
					FmtStatus("POST /TestResources without required extension", resp.StatusCode, 400))
			},
		}},
	},
}
