package spec

import (
	"fmt"
	"strings"
)

var l7644_3_11 = rfcLines(rfc7644Sections, "19-section-3.11.txt")

var rfc7644_3_11 = []Requirement{
	// Section 3.11

	{
		ID:      "RFC7644-3.11-L3683",
		Level:   Must,
		Summary: "/Me response Location header must be the permanent resource location",
		Source: Source{
			RFC:       7644,
			Section:   "3.11",
			StartLine: l7644_3_11(16),
			EndLine:   l7644_3_11(19),
			StartCol:  66,
		},
		Feature:  Me,
		Testable: true,
		Tests: []Test{{
			Name: "me_location_header",
			Fn: func(r *Run) {
				resp, err := r.Client.Get("/Me")
				r.RequireOK(err)

				if resp.StatusCode == 501 || resp.StatusCode == 404 {
					r.Skipf("GET /Me returned %d, server does not implement /Me", resp.StatusCode)
					return
				}

				if resp.StatusCode == 200 {
					loc := resp.Header.Get("Location")
					contentLoc := resp.Header.Get("Content-Location")
					metaLoc := GetString(resp.Body, "meta", "location")

					permLoc := loc
					if permLoc == "" {
						permLoc = contentLoc
					}
					if permLoc == "" {
						permLoc = metaLoc
					}

					r.Check(permLoc != "" && !strings.HasSuffix(permLoc, "/Me"),
						fmt.Sprintf("GET /Me: location = %q, must be permanent resource URI (not /Me)", permLoc))
				}
			},
		}},
	},
}
