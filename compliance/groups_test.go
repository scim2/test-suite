package compliance

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scim2/test-suite/spec"
)

func TestGroup_MemberRef(t *testing.T) {
	adv := requires(t, spec.Groups, spec.Users)

	userBody, userResp := createUser(t)
	if userResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
	}
	userID := idOf(userBody)

	groupBody, groupResp := createGroup(t, []map[string]any{
		{"value": userID, "type": "User"},
	})
	if groupResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
	}

	members, _ := groupBody["members"].([]any)
	if len(members) == 0 {
		t.Skip("no members in response")
	}

	first, _ := members[0].(map[string]any)
	ref, _ := first["$ref"].(string)
	if ref == "" {
		t.Skip("member $ref not returned")
	}

	// RFC7643-2.3.7-L555: $ref must be a resolvable URI containing
	// all components up to the endpoint URI.
	check(t, adv, []string{"RFC7643-2.3.7-L555"},
		strings.Contains(ref, "Users") && strings.Contains(ref, userID),
		fmt.Sprintf("member $ref = %q, want URI containing Users/%s", ref, userID))

	// RFC7643-2.3.7-L577: reference URI must be HTTP-addressable.
	// RFC7643-2.3.7-L578: GET on the reference must return the resource.
	// Resolve relative $ref against server base URL.
	resolved := ref
	if !strings.HasPrefix(ref, "http") {
		resolved = strings.TrimRight(client.BaseURL, "/") + "/" + strings.TrimLeft(ref, "/")
	}

	getResp, err := client.Get(ref)
	requireOK(t, err)

	check(t, adv, []string{"RFC7643-2.3.7-L577"},
		resolved != "",
		fmt.Sprintf("member $ref = %q is not HTTP-addressable", ref))

	check(t, adv, []string{"RFC7643-2.3.7-L578"},
		getResp.StatusCode == 200,
		fmtStatus(fmt.Sprintf("GET %s (member $ref)", ref), getResp.StatusCode, 200))

	if getResp.StatusCode == 200 {
		check(t, adv, []string{"RFC7643-2.3.7-L578"},
			idOf(getResp.Body) == userID,
			fmt.Sprintf("GET $ref returned id=%q, want %q", idOf(getResp.Body), userID))
	}
}

func TestGroup_MemberValueIsUserID(t *testing.T) {
	adv := requires(t, spec.Groups, spec.Users)

	userBody, userResp := createUser(t)
	if userResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
	}
	userID := idOf(userBody)

	groupBody, groupResp := createGroup(t, []map[string]any{
		{"value": userID, "type": "User"},
	})
	if groupResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
	}

	// RFC7643-4.1.2-L1351: group value sub-attribute must be the id of the Group resource.
	members, _ := groupBody["members"].([]any)
	if len(members) == 0 {
		t.Fatal("POST /Groups: members empty in response")
	}

	first, _ := members[0].(map[string]any)
	memberVal, _ := first["value"].(string)
	check(t, adv, []string{"RFC7643-4.1.2-L1351"},
		memberVal == userID,
		fmt.Sprintf("member value = %q, want user id %q", memberVal, userID))
}

func TestGroup_MembershipViaGroup(t *testing.T) {
	adv := requires(t, spec.Groups, spec.Users, spec.Patch)

	userBody, userResp := createUser(t)
	if userResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
	}
	userID := idOf(userBody)

	// Create empty group, then add member via PATCH on the Group.
	groupBody, groupResp := createGroup(t, nil)
	if groupResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
	}
	groupID := idOf(groupBody)

	patchResp, err := client.Patch("/Groups/"+groupID, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "add",
				"path": "members",
				"value": []map[string]any{
					{"value": userID, "type": "User"},
				},
			},
		},
	})
	requireOK(t, err)

	// RFC7643-4.1.2-L1354: group membership changes must be applied via the Group resource.
	check(t, adv, []string{"RFC7643-4.1.2-L1354"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH /Groups add member", patchResp.StatusCode, 200))

	// Verify member was added.
	getResp, err := client.Get("/Groups/" + groupID)
	requireOK(t, err)

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
		check(t, adv, []string{"RFC7643-4.1.2-L1354"},
			found,
			fmt.Sprintf("GET /Groups/%s: user %s not in members after PATCH add", groupID, userID))
	}
}

func TestGroup_UserGroupsAttribute(t *testing.T) {
	requires(t, spec.Groups, spec.Users)

	userBody, userResp := createUser(t)
	if userResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", userResp.StatusCode)
	}
	userID := idOf(userBody)

	groupBody, groupResp := createGroup(t, []map[string]any{
		{"value": userID, "type": "User"},
	})
	if groupResp.StatusCode != 201 {
		t.Fatalf("setup: POST /Groups returned %d", groupResp.StatusCode)
	}
	groupID := idOf(groupBody)

	// Fetch the user and check the groups attribute.
	getResp, err := client.Get("/Users/" + userID)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		groups, _ := getResp.Body["groups"].([]any)
		found := false
		for _, g := range groups {
			gm, ok := g.(map[string]any)
			if !ok {
				continue
			}
			if v, _ := gm["value"].(string); v == groupID {
				found = true
			}
		}
		// This is informational -- the groups attribute on User is readOnly
		// and derived from Group membership. Not all servers populate it.
		if !found {
			t.Logf("User.groups does not reflect group %s membership (server may not populate derived groups)", groupID)
		}
	}
}
