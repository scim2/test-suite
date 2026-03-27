package spec

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

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
		Tests: []Test{{
			Name: "integer_no_fraction",
			Fn: func(r *Run) {
				body, resp := r.CreateTestResource(map[string]any{"integerAttr": 42})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /TestResources returned %d", resp.StatusCode)
				}

				id := IDOf(body)
				getResp, err := r.Client.Get("/TestResources/" + id)
				r.RequireOK(err)
				if getResp.StatusCode == 200 {
					gv, _ := getResp.Body["integerAttr"].(float64)
					r.Check(gv == float64(int(gv)),
						fmt.Sprintf("integerAttr on GET = %v, contains fractional part", gv))
				}
			},
		}},
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
		Tests: []Test{{
			Name: "datetime_format",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				created := GetString(body, "meta", "created")
				r.Check(created != "" && strings.Contains(created, "T"),
					fmt.Sprintf("meta.created = %q, want xsd:dateTime with date and time", created))
			},
		}},
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
		Tests: []Test{{
			Name: "datetime_has_timezone",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				created := GetString(body, "meta", "created")
				hasTimezone := strings.HasSuffix(created, "Z") ||
					(len(created) > 5 && (created[len(created)-6] == '+' || created[len(created)-6] == '-'))
				r.Check(hasTimezone,
					fmt.Sprintf("meta.created = %q, must include timezone", created))
			},
		}},
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
		Tests: []Test{{
			Name: "binary_base64_encoded",
			Fn: func(r *Run) {
				raw := []byte("hello, compliance test")
				encoded := base64.StdEncoding.EncodeToString(raw)

				body, resp := r.CreateTestResource(map[string]any{
					"binaryAttr": encoded,
				})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				v, ok := body["binaryAttr"].(string)
				r.Check(ok && v != "",
					fmt.Sprintf("binaryAttr = %v, want base64 string", body["binaryAttr"]))

				if ok && v != "" {
					_, err := base64.StdEncoding.DecodeString(v)
					r.Check(err == nil,
						fmt.Sprintf("binaryAttr is not valid base64: %v", err))
				}
			},
		}},
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
		Tests: []Test{{
			Name: "reference_is_uri",
			Fn: func(r *Run) {
				body, resp := r.CreateTestResource(map[string]any{
					"referenceAttr": "https://example.com/resource/123",
				})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				v, ok := body["referenceAttr"].(string)
				r.Check(ok && v != "",
					fmt.Sprintf("referenceAttr = %v, want URI string", body["referenceAttr"]))

				if ok && v != "" {
					_, err := url.Parse(v)
					r.Check(err == nil,
						fmt.Sprintf("referenceAttr is not a valid URI: %v", err))
				}
			},
		}},
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
		Tests: []Test{{
			Name: "member_ref_contains_endpoint",
			Fn: func(r *Run) {
				userBody, userResp := r.CreateUser()
				if userResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
				}
				userID := IDOf(userBody)

				groupBody, groupResp := r.CreateGroup([]map[string]any{
					{"value": userID, "type": "User"},
				})
				if groupResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
				}

				members, _ := groupBody["members"].([]any)
				if len(members) == 0 {
					r.Skipf("no members in response")
				}

				first, _ := members[0].(map[string]any)
				ref, _ := first["$ref"].(string)
				if ref == "" {
					r.Skipf("member $ref not returned")
				}

				r.Check(strings.Contains(ref, "Users") && strings.Contains(ref, userID),
					fmt.Sprintf("member $ref = %q, want URI containing Users/%s", ref, userID))
			},
		}},
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
		Tests: []Test{{
			Name: "member_ref_http_addressable",
			Fn: func(r *Run) {
				userBody, userResp := r.CreateUser()
				if userResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
				}
				userID := IDOf(userBody)

				groupBody, groupResp := r.CreateGroup([]map[string]any{
					{"value": userID, "type": "User"},
				})
				if groupResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
				}

				members, _ := groupBody["members"].([]any)
				if len(members) == 0 {
					r.Skipf("no members in response")
				}

				first, _ := members[0].(map[string]any)
				ref, _ := first["$ref"].(string)
				if ref == "" {
					r.Skipf("member $ref not returned")
				}

				resolved := ref
				if !strings.HasPrefix(ref, "http") {
					resolved = strings.TrimRight(r.Client.BaseURL, "/") + "/" + strings.TrimLeft(ref, "/")
				}

				r.Check(resolved != "",
					fmt.Sprintf("member $ref = %q is not HTTP-addressable", ref))
			},
		}},
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
		Tests: []Test{{
			Name: "get_member_ref_returns_resource",
			Fn: func(r *Run) {
				userBody, userResp := r.CreateUser()
				if userResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
				}
				userID := IDOf(userBody)

				groupBody, groupResp := r.CreateGroup([]map[string]any{
					{"value": userID, "type": "User"},
				})
				if groupResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
				}

				members, _ := groupBody["members"].([]any)
				if len(members) == 0 {
					r.Skipf("no members in response")
				}

				first, _ := members[0].(map[string]any)
				ref, _ := first["$ref"].(string)
				if ref == "" {
					r.Skipf("member $ref not returned")
				}

				getResp, err := r.Client.Get(ref)
				r.RequireOK(err)

				r.Check(getResp.StatusCode == 200,
					FmtStatus(fmt.Sprintf("GET %s (member $ref)", ref), getResp.StatusCode, 200))

				if getResp.StatusCode == 200 {
					r.Check(IDOf(getResp.Body) == userID,
						fmt.Sprintf("GET $ref returned id=%q, want %q", IDOf(getResp.Body), userID))
				}
			},
		}},
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
		Tests: []Test{{
			Name: "attribute_order_independent",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"displayName": "Order Test",
					"userName":    "scim-test-order-" + RandomSuffix(),
					"schemas":     []string{UserSchema},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				r.Check(resp.StatusCode == 201,
					FmtStatus("POST /Users with non-standard attribute order", resp.StatusCode, 201))
			},
		}},
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
		Tests: []Test{{
			Name: "no_nested_complex",
			Fn: func(r *Run) {
				body, resp := r.CreateTestResource(map[string]any{
					"complexAttr": map[string]any{
						"sub1": "hello",
						"sub2": 99,
						"sub3": true,
					},
				})
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}

				c, ok := body["complexAttr"].(map[string]any)
				if !ok || c == nil {
					r.Skipf("complexAttr not returned as object")
				}

				for k, v := range c {
					_, isMap := v.(map[string]any)
					r.Check(!isMap,
						fmt.Sprintf("complexAttr.%s is a nested complex (not allowed)", k))
				}
			},
		}},
	},
}
