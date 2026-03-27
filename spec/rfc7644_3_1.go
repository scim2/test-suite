package spec

import (
	"fmt"
	"strings"
)

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
		Tests: []Test{{
			Name: "message_schema_uri_prefix",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Users/nonexistent-" + RandomSuffix())
				r.RequireOK(err)

				if resp.StatusCode == 404 && resp.Body != nil {
					schemas := GetStringSlice(resp.Body, "schemas")
					for _, s := range schemas {
						r.Check(strings.HasPrefix(s, "urn:ietf:params:scim:api:messages:"),
							fmt.Sprintf("error schema %q does not have required prefix", s))
					}
				}
			},
		}},
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
		Tests: []Test{{
			Name: "message_schemas_not_discoverable",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Schemas")
				r.RequireOK(err)

				if resp.StatusCode != 200 || resp.Body == nil {
					r.Fatalf("GET /Schemas returned %d", resp.StatusCode)
				}

				resources, _ := resp.Body["Resources"].([]any)
				for _, res := range resources {
					s, ok := res.(map[string]any)
					if !ok {
						continue
					}
					id, _ := s["id"].(string)
					r.Check(!strings.HasPrefix(id, "urn:ietf:params:scim:api:messages:"),
						fmt.Sprintf("/Schemas contains message schema %q which should not be discoverable", id))
				}
			},
		}},
	},
}
