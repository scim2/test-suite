package spec

import "fmt"

var l7643_4_1 = rfcLines(rfc7643Sections, "15-section-4.1.txt")

var rfc7643_4_1 = []Requirement{
	// Section 4.1.1

	{
		ID:      "RFC7643-4.1.1-L1035",
		Level:   Must,
		Summary: "Each User must include a non-empty userName value",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(15),
			EndLine:   l7643_4_1(18),
			StartCol:  51,
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "username_non_empty",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				userName, _ := body["userName"].(string)
				r.Check(userName != "",
					"POST /Users: userName must be non-empty in response")
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.1-L1186",
		Level:   ShallNot,
		Summary: "Cleartext or hashed password values shall not be returnable",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(165),
			EndLine:   l7643_4_1(167),
			StartCol:  59,
			EndCol:    22,
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "password_not_returned",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-pw-" + RandomSuffix(),
					"password": "s3cret!",
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				_, hasPassword := resp.Body["password"]
				r.Check(!hasPassword,
					"POST /Users with password: password value was returned in response")
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.1-L1222",
		Level:   MustNot,
		Summary: "Password value must not be returned in any form (returned is never)",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.1",
			StartLine: l7643_4_1(201),
			EndLine:   l7643_4_1(204),
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "password_never_in_response",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
				}

				_, hasPassword := body["password"]
				r.Check(!hasPassword,
					"POST /Users: password must not be returned")
			},
		}},
	},

	// Section 4.1.2

	{
		ID:      "RFC7643-4.1.2-L1244",
		Level:   Should,
		Summary: "Email values should be specified per RFC 5321",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(223),
			EndLine:   l7643_4_1(227),
			EndCol:    36,
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "email_rfc5321",
			Fn: func(r *Run) {
				body, resp := r.CreateUser()
				if resp.StatusCode != 201 {
					r.Fatalf("setup: POST returned %d", resp.StatusCode)
				}
				id := IDOf(body)
				userName, _ := body["userName"].(string)

				putResp, err := r.Client.Put("/Users/"+id, map[string]any{
					"schemas":  []string{UserSchema},
					"userName": userName,
					"emails": []map[string]any{
						{"value": "valid@example.com", "type": "work"},
					},
				})
				r.RequireOK(err)

				if putResp.StatusCode == 200 && putResp.Body != nil {
					emails, _ := putResp.Body["emails"].([]any)
					for _, e := range emails {
						em, ok := e.(map[string]any)
						if !ok {
							continue
						}
						v, _ := em["value"].(string)
						atCount := 0
						for _, c := range v {
							if c == '@' {
								atCount++
							}
						}
						r.Check(atCount == 1 && len(v) > 3,
							fmt.Sprintf("email %q does not look like RFC 5321 format", v))
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.2-L1281",
		Level:   Must,
		Summary: "Photo URI resource must be an image file, not a web page",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(260),
			EndLine:   l7643_4_1(263),
			EndCol:    15,
		},
		Feature:  Users,
		Testable: false,
	},
	{
		ID:      "RFC7643-4.1.2-L1322",
		Level:   Must,
		Summary: "Country value must be in ISO 3166-1 alpha-2 code format",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(300),
			EndLine:   l7643_4_1(303),
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "country_iso3166",
			Fn: func(r *Run) {
				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-country-" + RandomSuffix(),
					"addresses": []map[string]any{
						{"type": "work", "country": "US"},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode == 201 && resp.Body != nil {
					addrs, _ := resp.Body["addresses"].([]any)
					for _, a := range addrs {
						addr, ok := a.(map[string]any)
						if !ok {
							continue
						}
						country, _ := addr["country"].(string)
						if country == "" {
							continue
						}
						r.Check(len(country) == 2 && country[0] >= 'A' && country[0] <= 'Z' && country[1] >= 'A' && country[1] <= 'Z',
							fmt.Sprintf("country = %q, want 2-letter uppercase ISO 3166-1 code", country))
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.2-L1351",
		Level:   Must,
		Summary: "Group value sub-attribute must be the id of the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(321),
			EndLine:   l7643_4_1(336),
			StartCol:  63,
		},
		Feature:  Groups,
		Testable: true,
		Tests: []Test{{
			Name: "member_value_is_user_id",
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
					r.Fatalf("POST /Groups: members empty in response")
				}

				first, _ := members[0].(map[string]any)
				memberVal, _ := first["value"].(string)
				r.Check(memberVal == userID,
					fmt.Sprintf("member value = %q, want user id %q", memberVal, userID))
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.2-L1354",
		Level:   Must,
		Summary: "Group membership changes must be applied via the Group resource",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(333),
			EndLine:   l7643_4_1(336),
			StartCol:  22,
		},
		Feature:  Groups,
		Testable: true,
		Tests: []Test{{
			Name: "membership_via_group",
			Fn: func(r *Run) {
				userBody, userResp := r.CreateUser()
				if userResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
				}
				userID := IDOf(userBody)

				groupBody, groupResp := r.CreateGroup(nil)
				if groupResp.StatusCode != 201 {
					r.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
				}
				groupID := IDOf(groupBody)

				patchResp, err := r.Client.Patch("/Groups/"+groupID, map[string]any{
					"schemas": []string{PatchOpSchema},
					"Operations": []map[string]any{{
						"op":   "add",
						"path": "members",
						"value": []map[string]any{
							{"value": userID, "type": "User"},
						},
					}},
				})
				r.RequireOK(err)

				r.Check(patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
					FmtStatus("PATCH /Groups add member", patchResp.StatusCode, 200))

				getResp, err := r.Client.Get("/Groups/" + groupID)
				r.RequireOK(err)

				if getResp.StatusCode == 200 {
					members, _ := getResp.Body["members"].([]any)
					found := false
					for _, m := range members {
						mm, ok := m.(map[string]any)
						if !ok {
							continue
						}
						if v, _ := mm["value"].(string); v == userID {
							found = true
						}
					}
					r.Check(found,
						fmt.Sprintf("GET /Groups/%s: user %s not in members after PATCH add", groupID, userID))
				}
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.2-L1378",
		Level:   Must,
		Summary: "X.509 certificate values must be base64 encoded DER",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(355),
			EndLine:   l7643_4_1(359),
			EndCol:    41,
		},
		Feature:  Users,
		Testable: true,
		Tests: []Test{{
			Name: "x509_base64_der",
			Fn: func(r *Run) {
				certB64 := r.GenerateSelfSignedCertB64()

				resp, err := r.Client.Post("/Users", map[string]any{
					"schemas":  []string{UserSchema},
					"userName": "scim-test-x509-" + RandomSuffix(),
					"x509Certificates": []map[string]any{
						{"value": certB64},
					},
				})
				r.RequireOK(err)
				if resp.StatusCode == 201 {
					if id := IDOf(resp.Body); id != "" {
						r.Cleanup(func() { _, _ = r.Client.Delete("/Users/" + id) })
					}
				}

				if resp.StatusCode == 201 && resp.Body != nil {
					certs, _ := resp.Body["x509Certificates"].([]any)
					for _, c := range certs {
						cert, ok := c.(map[string]any)
						if !ok {
							continue
						}
						v, _ := cert["value"].(string)
						r.Check(v != "" && IsBase64(v),
							fmt.Sprintf("x509Certificate value is not valid base64: %q", v))
					}
				}
			},
		}},
	},
	{
		ID:      "RFC7643-4.1.2-L1379",
		Level:   MustNot,
		Summary: "A single certificate value must not contain multiple certificates",
		Source: Source{
			RFC:       7643,
			Section:   "4.1.2",
			StartLine: l7643_4_1(359),
			EndLine:   l7643_4_1(361),
			StartCol:  44,
		},
		Feature:  Users,
		Testable: false,
	},
}
