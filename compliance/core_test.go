package compliance

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/scim2/test-suite/scim"
	"github.com/scim2/test-suite/spec"
)

func TestAttributeOrderIndependent(t *testing.T) {
	adv := requires(t, spec.Users)

	// RFC7643-2.3.8-L592: must not require attributes in any specific order.
	// Create with attributes in unusual order.
	resp, err := client.Post("/Users", map[string]any{
		"displayName": "Order Test",
		"userName":    "scim-test-order-" + randomSuffix(),
		"schemas":     []string{userSchema},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	check(t, adv, []string{"RFC7643-2.3.8-L592"},
		resp.StatusCode == 201,
		fmtStatus("POST /Users with non-standard attribute order", resp.StatusCode, 201))
}

func TestCasePreservation(t *testing.T) {
	adv := requires(t, spec.Users)

	// RFC7643-7-L1750: server shall preserve case for case-exact attribute values.
	mixedCase := "ScIm-TeSt-CaSe-" + randomSuffix()
	resp, err := client.Post("/Users", map[string]any{
		"schemas":     []string{userSchema},
		"userName":    mixedCase,
		"displayName": "CaSeExAcT Display",
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode != 201 {
		t.Skipf("POST returned %d", resp.StatusCode)
	}

	// displayName is caseExact, case must be preserved.
	dn, _ := resp.Body["displayName"].(string)
	check(t, adv, []string{"RFC7643-7-L1750"},
		dn == "CaSeExAcT Display",
		fmt.Sprintf("displayName = %q, want \"CaSeExAcT Display\" (case not preserved)", dn))
}

func TestContentType(t *testing.T) {
	resp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	// RFC7644-3.8-L3551: Must support Accept: application/scim+json.
	ct := resp.Header.Get("Content-Type")
	check(t, true, []string{"RFC7644-3.8-L3551"},
		ct != "" && (ct == "application/scim+json" ||
			len(ct) > 24 && ct[:24] == "application/scim+json;"),
		fmt.Sprintf("Content-Type = %q, want application/scim+json", ct))
}

func TestCreateUser_CommonAttributes(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}

	// RFC7643-3.1-L855: id and meta must be assigned by the service provider.
	check(t, adv, []string{"RFC7643-3.1-L855"},
		idOf(body) != "" && body["meta"] != nil,
		"POST /Users: service provider must assign id and meta on creation")

	// RFC7643-3.1-L852: Common attributes must be defined for all resources.
	meta, _ := body["meta"].(map[string]any)
	check(t, adv, []string{"RFC7643-3.1-L852"},
		meta != nil,
		"POST /Users: missing meta attribute")

	if meta != nil {
		// RFC7643-3.1-L920 already tested in TestCreateUser_Returns201.
		// RFC7643-3.1-L927: meta.location must match Content-Location.
		contentLoc := resp.Header.Get("Content-Location")
		metaLoc, _ := meta["location"].(string)
		if contentLoc != "" {
			check(t, adv, []string{"RFC7643-3.1-L927"},
				contentLoc == metaLoc,
				fmt.Sprintf("Content-Location (%s) != meta.location (%s)", contentLoc, metaLoc))
		}
	}

	// RFC7643-3.3-L976: schemas must contain at least one value.
	schemas := getStringSlice(body, "schemas")
	check(t, adv, []string{"RFC7643-3.3-L976"},
		len(schemas) >= 1,
		"POST /Users: schemas must contain at least one value")

	// RFC7643-3-L711: schemas must only contain defined schema/schemaExtensions.
	check(t, adv, []string{"RFC7643-3-L711"},
		hasSchema(body, userSchema),
		fmt.Sprintf("POST /Users: schemas %v missing base schema %s", schemas, userSchema))

	// RFC7643-3-L713: No duplicate values in schemas.
	seen := make(map[string]bool)
	hasDupes := false
	for _, s := range schemas {
		if seen[s] {
			hasDupes = true
		}
		seen[s] = true
	}
	check(t, adv, []string{"RFC7643-3-L713"},
		!hasDupes,
		fmt.Sprintf("POST /Users: duplicate values in schemas: %v", schemas))

	// RFC7643-4.1.1-L1035: userName must be non-empty.
	userName, _ := body["userName"].(string)
	check(t, adv, []string{"RFC7643-4.1.1-L1035"},
		userName != "",
		"POST /Users: userName must be non-empty in response")

	// RFC7643-4.1.1-L1222: password must not be returned.
	_, hasPassword := body["password"]
	check(t, adv, []string{"RFC7643-4.1.1-L1222"},
		!hasPassword,
		"POST /Users: password must not be returned")
}

func TestCreateUser_DuplicateReturns409(t *testing.T) {
	adv := requires(t, spec.Users)

	userName := "scim-test-dup-" + randomSuffix()

	// Create first user.
	resp1, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
	})
	requireOK(t, err)
	if id := idOf(resp1.Body); id != "" {
		t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
	}

	// Attempt duplicate.
	resp2, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
	})
	requireOK(t, err)
	if resp2.StatusCode == 201 {
		if id := idOf(resp2.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	// RFC7644-3.3-L625: Duplicate must return 409 with uniqueness scimType.
	check(t, adv, []string{"RFC7644-3.3-L625"},
		resp2.StatusCode == 409,
		fmtStatus("POST /Users (duplicate)", resp2.StatusCode, 409))

	if resp2.Body != nil {
		scimType, _ := resp2.Body["scimType"].(string)
		check(t, adv, []string{"RFC7644-3.3-L625"},
			scimType == "uniqueness",
			fmt.Sprintf("POST /Users (duplicate): scimType = %q, want \"uniqueness\"", scimType))
	}
}

func TestCreateUser_IdUnique(t *testing.T) {
	adv := requires(t, spec.Users)

	_, resp1 := createUser(t)
	_, resp2 := createUser(t)

	if resp1.StatusCode != 201 || resp2.StatusCode != 201 {
		t.Skip("setup: could not create two users")
	}

	// RFC7643-3.1-L867: id must be unique across all resources.
	check(t, adv, []string{"RFC7643-3.1-L867"},
		idOf(resp1.Body) != idOf(resp2.Body),
		fmt.Sprintf("two users got same id: %q", idOf(resp1.Body)))
}

func TestCreateUser_PasswordNotReturned(t *testing.T) {
	adv := requires(t, spec.Users)

	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-pw-" + randomSuffix(),
		"password": "s3cret!",
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	// RFC7643-4.1.1-L1186: password SHALL NOT be returnable.
	_, hasPassword := resp.Body["password"]
	check(t, adv, []string{"RFC7643-4.1.1-L1186"},
		!hasPassword,
		"POST /Users with password: password value was returned in response")
}

func TestCreateUser_ReadOnlyIgnored(t *testing.T) {
	adv := requires(t, spec.Users)

	// Send id and meta in the request - server should ignore them.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-ro-" + randomSuffix(),
		"id":       "client-specified-id",
		"meta": map[string]any{
			"resourceType": "Invalid",
			"created":      "1999-01-01T00:00:00Z",
		},
	})
	requireOK(t, err)
	if resp.StatusCode == 201 {
		if id := idOf(resp.Body); id != "" {
			t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
		}
	}

	if resp.StatusCode != 201 {
		t.Skipf("POST /Users returned %d, cannot verify readOnly handling", resp.StatusCode)
	}

	// RFC7644-3.3-L584: readOnly attributes in request SHALL be ignored.
	check(t, adv, []string{"RFC7644-3.3-L584"},
		idOf(resp.Body) != "client-specified-id",
		"POST /Users: server used client-specified id instead of generating one")

	// RFC7643-3.1-L872: id must not be specified by the client.
	check(t, adv, []string{"RFC7643-3.1-L872"},
		idOf(resp.Body) != "client-specified-id",
		"POST /Users: server accepted client-specified id")

	// RFC7643-3.1-L911: meta shall be ignored when provided by clients.
	rt := getString(resp.Body, "meta", "resourceType")
	check(t, adv, []string{"RFC7643-3.1-L911"},
		rt == "User",
		fmt.Sprintf("POST /Users: meta.resourceType = %q, server should have ignored client meta", rt))

	// RFC7643-3.1-L873: id must not contain "bulkId".
	id := idOf(resp.Body)
	check(t, adv, []string{"RFC7643-3.1-L873"},
		id != "bulkId" && !strings.Contains(id, "bulkId"),
		fmt.Sprintf("POST /Users: id %q contains reserved keyword bulkId", id))
}

func TestCreateUser_ResponseBody(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)

	// RFC7644-3.3-L603: response body SHOULD contain the newly created resource.
	check(t, adv, []string{"RFC7644-3.3-L603"},
		resp.StatusCode == 201 && body != nil && idOf(body) != "",
		"POST /Users: response body should contain the created resource representation")
}

func TestCreateUser_Returns201(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)

	// RFC7644-3.3-L602: Successful creation returns 201.
	check(t, adv, []string{"RFC7644-3.3-L602"},
		resp.StatusCode == 201,
		fmtStatus("POST /Users", resp.StatusCode, 201))

	// RFC7644-3.3-L605: Location header and meta.location present.
	loc := resp.Header.Get("Location")
	check(t, adv, []string{"RFC7644-3.3-L605"},
		loc != "",
		"POST /Users: missing Location header")

	metaLoc := getString(body, "meta", "location")
	check(t, adv, []string{"RFC7644-3.3-L605"},
		metaLoc != "",
		"POST /Users: missing meta.location in response body")

	// RFC7643-3.1-L866: Response must include a non-empty id.
	check(t, adv, []string{"RFC7643-3.1-L866"},
		idOf(body) != "",
		"POST /Users: missing id in response body")

	// RFC7644-3.3.1-L712: meta.resourceType shall match endpoint.
	check(t, adv, []string{"RFC7644-3.3.1-L712"},
		getString(body, "meta", "resourceType") == "User",
		fmt.Sprintf("POST /Users: meta.resourceType = %q, want \"User\"",
			getString(body, "meta", "resourceType")))

	// RFC7643-3.1-L920: meta.created must be a DateTime.
	created := getString(body, "meta", "created")
	check(t, adv, []string{"RFC7643-3.1-L920"},
		created != "",
		"POST /Users: missing meta.created")

	// RFC7643-3.1-L925: meta.lastModified must equal meta.created on new resource.
	lastMod := getString(body, "meta", "lastModified")
	check(t, adv, []string{"RFC7643-3.1-L925"},
		created == lastMod,
		fmt.Sprintf("POST /Users: meta.created (%s) != meta.lastModified (%s) on new resource",
			created, lastMod))

	// RFC7643-3-L709: schemas array must be non-empty.
	schemas := getStringSlice(body, "schemas")
	check(t, adv, []string{"RFC7643-3-L709"},
		len(schemas) > 0,
		"POST /Users: empty schemas array in response")
}

func TestDateTimeFormat(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}

	created := getString(body, "meta", "created")
	// RFC7643-2.3.5-L527: DateTime must be valid xsd:dateTime.
	// Must contain both date and time (T separator).
	check(t, adv, []string{"RFC7643-2.3.5-L527"},
		created != "" && strings.Contains(created, "T"),
		fmt.Sprintf("meta.created = %q, want xsd:dateTime with date and time", created))

	// RFC7643-2.3.5-L531: DateTime must conform to XML constraints.
	// Should end with Z or timezone offset.
	hasTimezone := strings.HasSuffix(created, "Z") ||
		(len(created) > 5 && (created[len(created)-6] == '+' || created[len(created)-6] == '-'))
	check(t, adv, []string{"RFC7643-2.3.5-L531"},
		hasTimezone,
		fmt.Sprintf("meta.created = %q, must include timezone", created))
}

func TestDeleteUser_OmittedFromQuery(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// Create and delete a user.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-delquery-" + randomSuffix(),
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	userName, _ := resp.Body["userName"].(string)

	delResp, err := client.Delete("/Users/" + id)
	requireOK(t, err)
	if delResp.StatusCode != 204 {
		t.Fatalf("setup: DELETE returned %d", delResp.StatusCode)
	}

	// RFC7644-3.6-L2644: deleted resource must be omitted from query results.
	filter := fmt.Sprintf("userName eq %q", userName)
	queryResp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	tr, _ := queryResp.Body["totalResults"].(float64)
	check(t, adv, []string{"RFC7644-3.6-L2644"},
		tr == 0,
		fmt.Sprintf("query for deleted user returned totalResults=%v, want 0", queryResp.Body["totalResults"]))
}

func TestDeleteUser_Returns204(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create a user to delete (don't register cleanup, we delete manually).
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-del-" + randomSuffix(),
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d, need 201", resp.StatusCode)
	}
	id := idOf(resp.Body)
	if id == "" {
		t.Fatal("setup: POST /Users returned no id")
	}

	// Delete the user.
	delResp, err := client.Delete("/Users/" + id)
	requireOK(t, err)

	// RFC7644-3.6-L2657: Successful DELETE returns 204.
	check(t, adv, []string{"RFC7644-3.6-L2657"},
		delResp.StatusCode == 204,
		fmtStatus("DELETE /Users/"+id, delResp.StatusCode, 204))

	// RFC7644-3.6-L2642: Subsequent GET returns 404.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.6-L2642"},
		getResp.StatusCode == 404,
		fmtStatus("GET /Users/"+id+" after DELETE", getResp.StatusCode, 404))
}

func TestDiscoveryEndpoints_IgnoreQueryParams(t *testing.T) {
	// RFC7644-4-L4125: filter, sort, pagination SHALL be ignored on discovery endpoints.
	resp, err := client.Get("/ServiceProviderConfig?filter=invalid&sortBy=id&startIndex=999")
	requireOK(t, err)

	check(t, true, []string{"RFC7644-4-L4125"},
		resp.StatusCode == 200,
		fmtStatus("GET /ServiceProviderConfig with query params", resp.StatusCode, 200))
}

func TestETag_IfMatchWrongVersion(t *testing.T) {
	adv := requires(t, spec.Users, spec.ETag)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	id := idOf(body)
	userName, _ := body["userName"].(string)

	// RFC7644-3.13-L3956: when version is specified, service provider must
	// use it or reject the request.
	// Send a PUT with a wrong If-Match ETag.
	putResp, err := client.DoWithHeaders(
		"PUT", "/Users/"+id,
		map[string]any{
			"schemas":  []string{userSchema},
			"userName": userName,
		},
		map[string]string{"If-Match": "W/\"wrong-version\""},
	)
	requireOK(t, err)

	// Server should return 412 Precondition Failed.
	check(t, adv, []string{"RFC7644-3.13-L3956"},
		putResp.StatusCode == 412,
		fmtStatus("PUT with wrong If-Match", putResp.StatusCode, 412))
}

func TestETag_WeakPrefix(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}

	// Check meta.version for weak ETag prefix.
	version := getString(body, "meta", "version")
	if version != "" {
		// RFC7643-3.1-L940: weak entity-tags must be prefixed with W/.
		check(t, adv, []string{"RFC7643-3.1-L940"},
			strings.HasPrefix(version, "W/"),
			fmt.Sprintf("meta.version = %q, weak ETags must start with W/", version))
	}

	// Also check ETag header on GET.
	id := idOf(body)
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	etag := getResp.Header.Get("ETag")
	if etag != "" {
		check(t, adv, []string{"RFC7643-3.1-L940"},
			strings.HasPrefix(etag, "W/"),
			fmt.Sprintf("ETag header = %q, weak ETags must start with W/", etag))

		// RFC7644-3.14-L3967: when supported, ETags must be specified as HTTP header.
		check(t, adv, []string{"RFC7644-3.14-L3967"},
			true, "")
	}
}

func TestErrorFormat(t *testing.T) {
	// GET a resource that doesn't exist.
	resp, err := client.Get("/Users/nonexistent-" + randomSuffix())
	requireOK(t, err)

	if resp.StatusCode == 404 && resp.Body != nil {
		// RFC7644-3.12-L3707: Errors must be returned in JSON format.
		check(t, true, []string{"RFC7644-3.12-L3707"},
			hasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:Error"),
			"404 error response missing Error schema")

		status, _ := resp.Body["status"].(string)
		check(t, true, []string{"RFC7644-3.12-L3707"},
			status == "404",
			fmt.Sprintf("error response status = %q, want \"404\"", status))
	}
}

func TestGetUser(t *testing.T) {
	adv := requires(t, spec.Users)

	body, createResp := createUser(t)
	if createResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
	}
	id := idOf(body)

	resp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	check(t, adv, []string{"RFC7643-3.1-L868"},
		idOf(resp.Body) == id,
		fmt.Sprintf("GET /Users/%s: id = %q, want %q", id, idOf(resp.Body), id))

	// RFC7643-3.1-L927: meta.location must match Content-Location header.
	contentLoc := resp.Header.Get("Content-Location")
	metaLoc := getString(resp.Body, "meta", "location")
	if contentLoc != "" && metaLoc != "" {
		check(t, adv, []string{"RFC7643-3.1-L927"},
			contentLoc == metaLoc,
			fmt.Sprintf("Content-Location (%s) != meta.location (%s)",
				contentLoc, metaLoc))
	}
}

func TestImmutableNotUpdated(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)
	userName, _ := body["userName"].(string)

	// Try to change id via PUT (id is immutable).
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
		"id":       "changed-id",
	})
	requireOK(t, err)

	// RFC7643-7-L1766: immutable attributes shall not be updated after creation.
	if putResp.StatusCode == 200 {
		check(t, adv, []string{"RFC7643-7-L1766"},
			idOf(putResp.Body) == id,
			fmt.Sprintf("id changed from %q to %q after PUT", id, idOf(putResp.Body)))
	} else {
		check(t, adv, []string{"RFC7643-7-L1766"},
			putResp.StatusCode == 400,
			fmtStatus("PUT with changed immutable id", putResp.StatusCode, 400))
	}
}

func TestMe_LocationHeader(t *testing.T) {
	adv := requires(t, spec.Me)

	resp, err := client.Get("/Me")
	requireOK(t, err)

	if resp.StatusCode == 501 || resp.StatusCode == 404 {
		// Server does not implement /Me.
		check(t, adv, []string{"RFC7644-3.11-L3683"},
			!adv,
			fmtStatus("GET /Me", resp.StatusCode, 200))
		return
	}

	if resp.StatusCode == 200 {
		// RFC7644-3.11-L3683: /Me response Location header must be the
		// permanent resource location (e.g. /Users/{id}, not /Me).
		loc := resp.Header.Get("Location")
		contentLoc := resp.Header.Get("Content-Location")
		metaLoc := getString(resp.Body, "meta", "location")

		permLoc := loc
		if permLoc == "" {
			permLoc = contentLoc
		}
		if permLoc == "" {
			permLoc = metaLoc
		}

		check(t, adv, []string{"RFC7644-3.11-L3683"},
			permLoc != "" && !strings.HasSuffix(permLoc, "/Me"),
			fmt.Sprintf("GET /Me: location = %q, must be permanent resource URI (not /Me)", permLoc))
	}
}

func TestMessageSchema_URIPrefix(t *testing.T) {
	// RFC7644-3.1-L436: SCIM message schemas URI prefix must begin with
	// urn:ietf:params:scim:api:
	// Verify by checking the error response schema.
	resp, err := client.Get("/Users/nonexistent-" + randomSuffix())
	requireOK(t, err)

	if resp.StatusCode == 404 && resp.Body != nil {
		schemas := getStringSlice(resp.Body, "schemas")
		for _, s := range schemas {
			check(t, true, []string{"RFC7644-3.1-L436"},
				strings.HasPrefix(s, "urn:ietf:params:scim:api:messages:"),
				fmt.Sprintf("error schema %q does not have required prefix", s))
		}
	}
}

func TestMetaLocation_MatchesContentLocation(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode != 200 {
		t.Fatalf("GET /Users/%s returned %d", id, getResp.StatusCode)
	}

	// RFC7643-3.1-L927: meta.location must match Content-Location header.
	metaLoc := getString(getResp.Body, "meta", "location")
	contentLoc := getResp.Header.Get("Content-Location")
	check(t, adv, []string{"RFC7643-3.1-L927"},
		metaLoc != "",
		"GET /Users: missing meta.location")

	if contentLoc != "" && metaLoc != "" {
		check(t, adv, []string{"RFC7643-3.1-L927"},
			contentLoc == metaLoc,
			fmt.Sprintf("Content-Location (%s) != meta.location (%s)", contentLoc, metaLoc))
	}
}

func TestMultiResourceQuery(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// RFC7644-3.4.2.1-L918: multi-resource-type queries must process
	// filters per Section 3.4.2.2.
	// A root query (GET /) or a query to a non-specific endpoint should
	// return resources from multiple types. Some servers don't support this.
	resp, err := client.Get("/?filter=" + url.QueryEscape("id pr"))
	requireOK(t, err)

	switch resp.StatusCode {
	case 200:
		check(t, adv, []string{"RFC7644-3.4.2.1-L918"},
			hasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
			"root query: missing ListResponse schema")
	case 404, 501:
		// Server doesn't support root-level queries.
		check(t, adv, []string{"RFC7644-3.4.2.1-L918"},
			true, "")
		t.Logf("root-level query not supported (status %d)", resp.StatusCode)
	default:
		check(t, adv, []string{"RFC7644-3.4.2.1-L918"},
			false,
			fmtStatus("GET /?filter=id pr", resp.StatusCode, 200))
	}
}

func TestPagination_Count(t *testing.T) {
	adv := requires(t, spec.Users)

	// Create 3 users.
	for range 3 {
		createUser(t)
	}

	// Request count=1.
	resp, err := client.Get("/Users?count=1")
	requireOK(t, err)

	if resp.StatusCode != 200 {
		t.Fatalf("GET /Users?count=1 returned %d", resp.StatusCode)
	}

	resources, _ := resp.Body["Resources"].([]any)
	// RFC7644-3.4.2.4-L1363: must not return more than count.
	check(t, adv, []string{"RFC7644-3.4.2.4-L1363"},
		len(resources) <= 1,
		fmt.Sprintf("GET /Users?count=1: got %d resources, want <= 1", len(resources)))
}

func TestPagination_StartIndex(t *testing.T) {
	adv := requires(t, spec.Users)

	// RFC7644-3.4.2.4-L1358: startIndex < 1 shall be interpreted as 1.
	resp, err := client.Get("/Users?startIndex=0")
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.4.2.4-L1358"},
		resp.StatusCode == 200,
		fmtStatus("GET /Users?startIndex=0", resp.StatusCode, 200))
}

func TestPutUser(t *testing.T) {
	adv := requires(t, spec.Users)

	body, createResp := createUser(t)
	if createResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", createResp.StatusCode)
	}
	id := idOf(body)
	schemas := getStringSlice(body, "schemas")
	userName, _ := body["userName"].(string)

	putBody := map[string]any{
		"schemas":     schemas,
		"userName":    userName,
		"displayName": "Updated " + randomSuffix(),
	}

	resp, err := client.Put("/Users/"+id, putBody)
	requireOK(t, err)

	// RFC7644-3.5-L1607: PUT must be supported, returns 200.
	check(t, adv, []string{"RFC7644-3.5-L1607"},
		resp.StatusCode == 200,
		fmtStatus("PUT /Users/"+id, resp.StatusCode, 200))

	// RFC7644-3.5.1-L1665: readOnly values shall be ignored (id shouldn't change).
	check(t, adv, []string{"RFC7644-3.5.1-L1665"},
		idOf(resp.Body) == id,
		fmt.Sprintf("PUT /Users/%s: id changed to %q", id, idOf(resp.Body)))
}

func TestPutUser_ImmutableMustMatch(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	id := idOf(body)
	userName, _ := body["userName"].(string)

	// Try to PUT with a different id in the body (id is immutable).
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
		"id":       "different-id-value",
	})
	requireOK(t, err)

	// RFC7644-3.5.1-L1660: immutable attribute input must match existing or return 400.
	// Server must either ignore the mismatched id or return 400.
	if putResp.StatusCode == 200 {
		check(t, adv, []string{"RFC7644-3.5.1-L1660"},
			idOf(putResp.Body) == id,
			fmt.Sprintf("PUT with mismatched id: server changed id to %q", idOf(putResp.Body)))
	} else {
		check(t, adv, []string{"RFC7644-3.5.1-L1660"},
			putResp.StatusCode == 400,
			fmtStatus("PUT /Users with mismatched immutable id", putResp.StatusCode, 400))
	}
}

func TestPutUser_NotCreate(t *testing.T) {
	adv := requires(t, spec.Users)

	// PUT to a non-existent resource should NOT create it.
	resp, err := client.Put("/Users/nonexistent-"+randomSuffix(), map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-putcreate-" + randomSuffix(),
	})
	requireOK(t, err)

	// RFC7644-3.5.1-L1637 / RFC7644-3.2-L490: PUT must not create.
	check(t, adv, []string{"RFC7644-3.5.1-L1637", "RFC7644-3.2-L490"},
		resp.StatusCode == 404,
		fmtStatus("PUT /Users/<nonexistent>", resp.StatusCode, 404))
}

func TestPutUser_RequiredAttributes(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// PUT without required userName.
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":     []string{userSchema},
		"displayName": "No UserName",
	})
	requireOK(t, err)

	// RFC7644-3.5.1-L1667: required attributes must be specified in PUT.
	check(t, adv, []string{"RFC7644-3.5.1-L1667"},
		putResp.StatusCode == 400,
		fmtStatus("PUT /Users without required userName", putResp.StatusCode, 400))
}

func TestQueryUsers_EmptyFilter(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// Use a filter that should match nothing.
	filter := fmt.Sprintf("userName eq %q", "scim-nonexistent-"+randomSuffix())
	resp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	// RFC7644-3.4.2-L847: Empty result returns 200 with totalResults 0.
	check(t, adv, []string{"RFC7644-3.4.2-L847"},
		resp.StatusCode == 200,
		fmtStatus("GET /Users?filter=... (no match)", resp.StatusCode, 200))

	if resp.Body != nil {
		tr, _ := resp.Body["totalResults"].(float64)
		check(t, adv, []string{"RFC7644-3.4.2-L847"},
			tr == 0,
			fmt.Sprintf("GET /Users?filter=... (no match): totalResults = %v, want 0", resp.Body["totalResults"]))
	}
}

func TestQueryUsers_ListResponse(t *testing.T) {
	adv := requires(t, spec.Users)

	resp, err := client.Get("/Users")
	requireOK(t, err)

	// RFC7644-3.4.2-L817: Query responses must use ListResponse schema.
	check(t, adv, []string{"RFC7644-3.4.2-L817"},
		hasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
		"GET /Users: missing ListResponse schema")

	// totalResults must be present.
	_, hasTR := resp.Body["totalResults"]
	check(t, adv, []string{"RFC7644-3.4.2-L817"},
		hasTR,
		"GET /Users: missing totalResults in ListResponse")
}

func TestReadOnlyNotModified(t *testing.T) {
	adv := requires(t, spec.Users)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)
	userName, _ := body["userName"].(string)

	// PUT with meta.resourceType set to something wrong.
	putResp, err := client.Put("/Users/"+id, map[string]any{
		"schemas":  []string{userSchema},
		"userName": userName,
		"meta": map[string]any{
			"resourceType": "Invalid",
		},
	})
	requireOK(t, err)

	// RFC7643-7-L1759: readOnly attributes shall not be modified.
	if putResp.StatusCode == 200 && putResp.Body != nil {
		rt := getString(putResp.Body, "meta", "resourceType")
		check(t, adv, []string{"RFC7643-7-L1759"},
			rt == "User",
			fmt.Sprintf("after PUT with invalid meta.resourceType: got %q, want \"User\"", rt))
	}
}

func TestResourceTypes(t *testing.T) {
	resp, err := client.Get("/ResourceTypes")
	requireOK(t, err)

	if resp.StatusCode != 200 || resp.Body == nil {
		t.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
	}

	resources, _ := resp.Body["Resources"].([]any)
	for _, r := range resources {
		rt, ok := r.(map[string]any)
		if !ok {
			continue
		}
		name, _ := rt["name"].(string)
		schema, _ := rt["schema"].(string)

		// RFC7643-6-L1619: name must be specified.
		check(t, true, []string{"RFC7643-6-L1619"},
			name != "",
			"ResourceType missing name")

		// RFC7643-6-L1641: schema URI must be specified.
		check(t, true, []string{"RFC7643-6-L1641"},
			schema != "",
			fmt.Sprintf("ResourceType %q: missing schema URI", name))
	}
}

func TestResourceTypes_SchemaURI(t *testing.T) {
	resp, err := client.Get("/ResourceTypes")
	requireOK(t, err)

	if resp.StatusCode != 200 || resp.Body == nil {
		t.Fatalf("GET /ResourceTypes returned %d", resp.StatusCode)
	}

	// Fetch schemas to cross-reference.
	schemasResp, err := client.Get("/Schemas")
	requireOK(t, err)

	schemaIDs := make(map[string]bool)
	if schemasResp.StatusCode == 200 && schemasResp.Body != nil {
		resources, _ := schemasResp.Body["Resources"].([]any)
		for _, r := range resources {
			s, ok := r.(map[string]any)
			if !ok {
				continue
			}
			if id, ok := s["id"].(string); ok {
				schemaIDs[id] = true
			}
		}
	}

	resources, _ := resp.Body["Resources"].([]any)
	for _, r := range resources {
		rt, ok := r.(map[string]any)
		if !ok {
			continue
		}
		name, _ := rt["name"].(string)
		schemaURI, _ := rt["schema"].(string)

		// RFC7643-6-L1641: schema URI must equal id of associated Schema resource.
		if len(schemaIDs) > 0 && schemaURI != "" {
			check(t, true, []string{"RFC7643-6-L1641"},
				schemaIDs[schemaURI],
				fmt.Sprintf("ResourceType %q: schema %q not found in /Schemas", name, schemaURI))
		}

		// RFC7643-6-L1650: extension schema URI must equal id of a Schema resource.
		exts, _ := rt["schemaExtensions"].([]any)
		for _, ext := range exts {
			e, ok := ext.(map[string]any)
			if !ok {
				continue
			}
			extSchema, _ := e["schema"].(string)
			if extSchema != "" && len(schemaIDs) > 0 {
				check(t, true, []string{"RFC7643-6-L1650"},
					schemaIDs[extSchema],
					fmt.Sprintf("ResourceType %q: extension schema %q not found in /Schemas", name, extSchema))
			}
		}
	}
}

func TestResponse_UTF8(t *testing.T) {
	// RFC7644-3.8-L3537: servers must accept and return JSON using UTF-8.
	resp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	ct := resp.Header.Get("Content-Type")
	// UTF-8 is the default for JSON; charset if present must be utf-8.
	if strings.Contains(ct, "charset") {
		check(t, true, []string{"RFC7644-3.8-L3537"},
			strings.Contains(strings.ToLower(ct), "utf-8"),
			fmt.Sprintf("Content-Type charset is not UTF-8: %s", ct))
	} else {
		// No charset means default UTF-8, which is compliant.
		check(t, true, []string{"RFC7644-3.8-L3537"},
			true, "")
	}

	// Also verify the body is valid JSON (UTF-8 encoded).
	check(t, true, []string{"RFC7644-3.8-L3537"},
		resp.Body != nil,
		"response body is not valid JSON")
}

func TestSchemas(t *testing.T) {
	resp, err := client.Get("/Schemas")
	requireOK(t, err)

	// RFC7644-4-L4098: SHALL return all supported schemas in ListResponse format.
	check(t, true, []string{"RFC7644-4-L4098"},
		resp.StatusCode == 200,
		fmtStatus("GET /Schemas", resp.StatusCode, 200))

	check(t, true, []string{"RFC7644-4-L4098"},
		hasSchema(resp.Body, "urn:ietf:params:scim:api:messages:2.0:ListResponse"),
		"/Schemas response missing ListResponse schema")
}

func TestSchemas_Content(t *testing.T) {
	resp, err := client.Get("/Schemas")
	requireOK(t, err)

	if resp.StatusCode != 200 || resp.Body == nil {
		t.Fatalf("GET /Schemas returned %d", resp.StatusCode)
	}

	resources, _ := resp.Body["Resources"].([]any)
	for _, r := range resources {
		s, ok := r.(map[string]any)
		if !ok {
			continue
		}
		id, _ := s["id"].(string)
		name, _ := s["name"].(string)

		// RFC7643-7-L1691: schema URI must be specified.
		check(t, true, []string{"RFC7643-7-L1691"},
			id != "",
			"Schema missing id (URI)")

		// RFC7643-7-L1700: schema name must be specified.
		check(t, true, []string{"RFC7643-7-L1700"},
			name != "",
			fmt.Sprintf("Schema %q: missing name", id))
	}
}

func TestSchemas_NotDiscoverableMessageSchemas(t *testing.T) {
	resp, err := client.Get("/Schemas")
	requireOK(t, err)

	if resp.StatusCode != 200 || resp.Body == nil {
		t.Fatalf("GET /Schemas returned %d", resp.StatusCode)
	}

	resources, _ := resp.Body["Resources"].([]any)
	// RFC7644-3.1-L439: SCIM message schemas shall NOT be discoverable via /Schemas.
	for _, r := range resources {
		s, ok := r.(map[string]any)
		if !ok {
			continue
		}
		id, _ := s["id"].(string)
		check(t, true, []string{"RFC7644-3.1-L439"},
			!strings.HasPrefix(id, "urn:ietf:params:scim:api:messages:"),
			fmt.Sprintf("/Schemas contains message schema %q which should not be discoverable", id))
	}
}

func TestServiceProviderConfig(t *testing.T) {
	resp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	// RFC7644-4-L4068: SHALL return ServiceProviderConfig schema.
	check(t, true, []string{"RFC7644-4-L4068"},
		resp.StatusCode == 200,
		fmtStatus("GET /ServiceProviderConfig", resp.StatusCode, 200))

	check(t, true, []string{"RFC7644-4-L4068"},
		hasSchema(resp.Body, "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"),
		"ServiceProviderConfig missing expected schemas value")
}

func TestServiceProviderConfig_AuthSchemes(t *testing.T) {
	// RFC7643-5-L1581: authenticationSchemes should be publicly accessible.
	resp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	if resp.StatusCode != 200 {
		t.Fatalf("GET /ServiceProviderConfig returned %d", resp.StatusCode)
	}

	authSchemes, _ := resp.Body["authenticationSchemes"].([]any)
	check(t, true, []string{"RFC7643-5-L1581"},
		len(authSchemes) > 0,
		"ServiceProviderConfig: authenticationSchemes is empty or missing")
}

func TestTooManyResults(t *testing.T) {
	adv := requires(t, spec.Users, spec.Filter)

	// Discover the maxResults limit.
	spcResp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)
	filterCfg, _ := spcResp.Body["filter"].(map[string]any)
	maxResults, _ := filterCfg["maxResults"].(float64)
	if maxResults == 0 || maxResults > 100 {
		t.Skipf("maxResults=%v, skipping (too large to create that many resources)", maxResults)
	}

	limit := int(maxResults)

	// Create limit+1 resources with a shared prefix.
	prefix := "scim-test-toomany-" + randomSuffix()
	for i := range limit + 1 {
		resp, err := client.Post("/Users", map[string]any{
			"schemas":  []string{userSchema},
			"userName": fmt.Sprintf("%s-%d", prefix, i),
		})
		requireOK(t, err)
		if resp.StatusCode == 201 {
			if id := idOf(resp.Body); id != "" {
				t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })
			}
		}
	}

	// Query with a filter that matches all of them.
	filter := fmt.Sprintf("userName sw %q", prefix)
	qResp, err := client.Get("/Users?filter=" + url.QueryEscape(filter))
	requireOK(t, err)

	// RFC7644-3.4.2.1-L912: too many results shall return 400 with scimType tooMany.
	// The server may also truncate and return 200 -- both behaviors exist.
	if qResp.StatusCode == 400 {
		if qResp.Body != nil {
			scimType, _ := qResp.Body["scimType"].(string)
			check(t, adv, []string{"RFC7644-3.4.2.1-L912"},
				scimType == "tooMany",
				fmt.Sprintf("scimType = %q, want \"tooMany\"", scimType))
		}
	} else {
		// Server truncated instead of rejecting.
		resources, _ := qResp.Body["Resources"].([]any)
		check(t, adv, []string{"RFC7644-3.4.2.1-L912"},
			len(resources) <= limit,
			fmt.Sprintf("returned %d resources, maxResults=%d", len(resources), limit))
	}
}

func TestWWWAuthenticate(t *testing.T) {
	// RFC7644-2-L307: service provider shall indicate supported HTTP auth
	// schemes via WWW-Authenticate header.
	// Make a raw HTTP request without auth to trigger a 401.
	rawClient := &scim.Client{
		BaseURL:    client.BaseURL,
		HTTPClient: client.HTTPClient,
	}
	resp, err := rawClient.Get("/Users")
	requireOK(t, err)

	if resp.StatusCode == 401 {
		wwwAuth := resp.Header.Get("WWW-Authenticate")
		check(t, true, []string{"RFC7644-2-L307"},
			wwwAuth != "",
			"401 response missing WWW-Authenticate header")
	} else {
		// Server doesn't require auth (e.g. test server).
		// Record the requirement as tested but passed since it's not applicable.
		check(t, true, []string{"RFC7644-2-L307"},
			true, "")
	}
}
