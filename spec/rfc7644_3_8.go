package spec

import (
	"fmt"
	"strings"
)

var l7644_3_8 = rfcLines(rfc7644Sections, "16-section-3.8.txt")

var rfc7644_3_8 = []Requirement{
	// Section 3.8

	{
		ID:      "RFC7644-3.8-L3537",
		Level:   Must,
		Summary: "Servers must accept and return JSON using UTF-8 encoding",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: l7644_3_8(3),
			EndLine:   l7644_3_8(5),
			EndCol:    19,
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "response_utf8",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				ct := resp.Header.Get("Content-Type")
				if strings.Contains(ct, "charset") {
					r.Check(strings.Contains(strings.ToLower(ct), "utf-8"),
						fmt.Sprintf("Content-Type charset is not UTF-8: %s", ct))
				} else {
					r.Check(true, "")
				}

				r.Check(resp.Body != nil,
					"response body is not valid JSON")
			},
		}},
	},
	{
		ID:      "RFC7644-3.8-L3551",
		Level:   Must,
		Summary: "Service providers must support Accept: application/scim+json",
		Source: Source{
			RFC:       7644,
			Section:   "3.8",
			StartLine: l7644_3_8(17),
			EndLine:   l7644_3_8(21),
		},
		Feature:  Core,
		Testable: true,
		Tests: []Test{{
			Name: "content_type_scim_json",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/ServiceProviderConfig")
				r.RequireOK(err)

				ct := resp.Header.Get("Content-Type")
				r.Check(ct != "" && (ct == "application/scim+json" ||
					len(ct) > 24 && ct[:24] == "application/scim+json;"),
					fmt.Sprintf("Content-Type = %q, want application/scim+json", ct))
			},
		}},
	},
}
