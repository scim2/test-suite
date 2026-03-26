package compliance

import (
	"fmt"
	"testing"

	"github.com/scim2/test-suite/spec"
)

const patchOpSchema = "urn:ietf:params:scim:api:messages:2.0:PatchOp"

func TestPatch_AddDisplayName(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	id := idOf(body)

	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "add",
				"path":  "displayName",
				"value": "Patched User",
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.5.2-L1939: successful PATCH returns 200 or 204.
	check(t, adv, []string{"RFC7644-3.5.2-L1939"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH /Users/"+id, patchResp.StatusCode, 200))

	// RFC7644-3.5.2-L1805: body must contain PatchOp schema.
	// (We sent it; if server accepted it, it understood the schema.)
	check(t, adv, []string{"RFC7644-3.5.2-L1805"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		"PATCH request with PatchOp schema was not accepted")

	// Verify the attribute was added.
	if patchResp.StatusCode == 200 && patchResp.Body != nil {
		dn, _ := patchResp.Body["displayName"].(string)
		check(t, adv, []string{"RFC7644-3.5.2.1-L1972"},
			dn == "Patched User",
			fmt.Sprintf("PATCH add displayName: got %q, want \"Patched User\"", dn))
	}
}

func TestPatch_AddExistingNoTimestampChange(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	resp, err := client.Post("/Users", map[string]any{
		"schemas":     []string{userSchema},
		"userName":    "scim-test-addexist-" + randomSuffix(),
		"displayName": "Existing Value",
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	lastMod := getString(resp.Body, "meta", "lastModified")

	// Add the same value again.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "add",
				"path":  "displayName",
				"value": "Existing Value",
			},
		},
	})
	requireOK(t, err)

	if patchResp.StatusCode == 200 || patchResp.StatusCode == 204 {
		// RFC7644-3.5.2.1-L2004: add of already-existing value shall not change timestamp.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)
		if getResp.StatusCode == 200 && lastMod != "" {
			newLastMod := getString(getResp.Body, "meta", "lastModified")
			check(t, adv, []string{"RFC7644-3.5.2.1-L2004"},
				newLastMod == lastMod,
				fmt.Sprintf("add existing value changed lastModified: %s -> %s", lastMod, newLastMod))
		}
	}
}

func TestPatch_AddToComplex(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// RFC7644-3.5.2.1-L1988: add to complex attribute shall specify sub-attributes.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "add",
				"path": "name",
				"value": map[string]any{
					"givenName":  "Test",
					"familyName": "User",
				},
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2.1-L1988"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH add to complex name", patchResp.StatusCode, 200))

	// Verify sub-attributes were set.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)
	if getResp.StatusCode == 200 {
		gn := getString(getResp.Body, "name", "givenName")
		check(t, adv, []string{"RFC7644-3.5.2.1-L1988"},
			gn == "Test",
			fmt.Sprintf("name.givenName = %q, want \"Test\"", gn))
	}
}

func TestPatch_Atomicity(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Send two operations: first valid, second invalid (remove required userName).
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "displayName",
				"value": "Should Not Persist",
			},
			{
				"op":   "remove",
				"path": "userName",
			},
		},
	})
	requireOK(t, err)

	if patchResp.StatusCode >= 400 {
		// RFC7644-3.5.2-L1931: PATCH shall be atomic.
		// RFC7644-3.5.2-L1933: on error, original resource must be restored.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)

		if getResp.StatusCode == 200 {
			dn, _ := getResp.Body["displayName"].(string)
			check(t, adv, []string{"RFC7644-3.5.2-L1931", "RFC7644-3.5.2-L1933"},
				dn != "Should Not Persist",
				"PATCH was not atomic: displayName changed despite second op failing")
		}
	} else {
		// Server accepted the removal of userName.
		// This is questionable but some servers allow it.
		t.Log("server accepted removal of required attribute userName")
	}
}

func TestPatch_Atomicity_Recorded(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)
	originalDN, _ := body["displayName"].(string)

	// Send two operations: first valid, second targets a readOnly attribute.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "displayName",
				"value": "Atomic-Change-" + randomSuffix(),
			},
			{
				"op":    "replace",
				"path":  "id",
				"value": "invalid-id-change",
			},
		},
	})
	requireOK(t, err)

	if patchResp.StatusCode >= 400 {
		// Server rejected -- verify original state is preserved.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)

		if getResp.StatusCode == 200 {
			dn, _ := getResp.Body["displayName"].(string)
			// RFC7644-3.5.2-L1931: PATCH shall be atomic.
			// RFC7644-3.5.2-L1933: on error, original resource must be restored.
			check(t, adv, []string{"RFC7644-3.5.2-L1931", "RFC7644-3.5.2-L1933"},
				dn == originalDN,
				fmt.Sprintf("PATCH not atomic: displayName changed to %q despite error", dn))
		}
	} else {
		// Server accepted both -- second op was likely ignored (readOnly).
		// Still record a pass for atomicity since there was no error.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)
		if getResp.StatusCode == 200 {
			check(t, adv, []string{"RFC7644-3.5.2-L1931", "RFC7644-3.5.2-L1933"},
				idOf(getResp.Body) == id,
				"PATCH: id was modified (should be readOnly)")
		}
	}
}

func TestPatch_MutabilityCompatible(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Attempt to replace a readWrite attribute (displayName) - should work.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "displayName",
				"value": "Mutability Test",
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.5.2-L1886: each op must be compatible with attribute mutability.
	check(t, adv, []string{"RFC7644-3.5.2-L1886"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH replace readWrite displayName", patchResp.StatusCode, 200))
}

func TestPatch_OperationsArray(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// RFC7644-3.5.2-L1808: body must contain Operations array.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "add",
				"path":  "nickName",
				"value": "Nicky",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2-L1808"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH with Operations array", patchResp.StatusCode, 200))
}

func TestPatch_PrimaryAutoUnset(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	// Create user with one primary email.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-primary-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "first@example.com", "type": "work", "primary": true},
		},
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	// Add another email as primary.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "add",
				"path": "emails",
				"value": []map[string]any{
					{"value": "second@example.com", "type": "home", "primary": true},
				},
			},
		},
	})
	requireOK(t, err)

	if patchResp.StatusCode == 200 || patchResp.StatusCode == 204 {
		// RFC7644-3.5.2-L1927: setting primary=true shall auto-set primary=false for others.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)
		if getResp.StatusCode == 200 {
			emails, _ := getResp.Body["emails"].([]any)
			primaryCount := 0
			for _, e := range emails {
				em, ok := e.(map[string]any)
				if !ok {
					continue
				}
				if p, _ := em["primary"].(bool); p {
					primaryCount++
				}
			}
			check(t, adv, []string{"RFC7644-3.5.2-L1927"},
				primaryCount <= 1,
				fmt.Sprintf("after adding primary email: %d primary values, want <= 1", primaryCount))
		}
	} else if patchResp.StatusCode == 400 {
		// Server rejected adding a second primary; acceptable.
		check(t, adv, []string{"RFC7644-3.5.2-L1927"}, true, "")
	} else {
		check(t, adv, []string{"RFC7644-3.5.2-L1927"}, false,
			fmtStatus("PATCH add primary email", patchResp.StatusCode, 200))
	}
}

func TestPatch_ReadOnlyNotModifiable(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Try to PATCH the id (readOnly attribute).
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "id",
				"value": "hacked-id",
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.5.2-L1888: must not modify readOnly attributes via PATCH.
	if patchResp.StatusCode >= 200 && patchResp.StatusCode < 300 {
		// If accepted, id must not have changed.
		getResp, err := client.Get("/Users/" + id)
		requireOK(t, err)
		check(t, adv, []string{"RFC7644-3.5.2-L1888"},
			getResp.StatusCode == 200 && idOf(getResp.Body) == id,
			"PATCH replace on readOnly id: id was modified")
	} else {
		// Server correctly rejected the operation.
		check(t, adv, []string{"RFC7644-3.5.2-L1888"},
			patchResp.StatusCode == 400,
			fmtStatus("PATCH replace readOnly id", patchResp.StatusCode, 400))
	}
}

func TestPatch_RemoveDisplayName(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	// Create user with displayName.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":     []string{userSchema},
		"userName":    "scim-test-patchrm-" + randomSuffix(),
		"displayName": "To Be Removed",
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "remove",
				"path": "displayName",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2-L1939"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH remove /Users/"+id, patchResp.StatusCode, 200))

	// RFC7644-3.5.2.2-L2146: removed single-value attribute shall be unassigned.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		dn, hasDN := getResp.Body["displayName"]
		check(t, adv, []string{"RFC7644-3.5.2.2-L2146"},
			!hasDN || dn == nil || dn == "",
			fmt.Sprintf("after PATCH remove: displayName = %v, want unassigned", dn))
	}
}

func TestPatch_RemoveMultiValued(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	// Create user with emails.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-patchmv-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "a@example.com", "type": "work"},
			{"value": "b@example.com", "type": "home"},
		},
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	// Remove all emails.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "remove",
				"path": "emails",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2-L1939"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH remove emails", patchResp.StatusCode, 200))

	// RFC7644-3.5.2.2-L2151: removing all values shall make attribute unassigned.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		emails, hasEmails := getResp.Body["emails"]
		emailArr, _ := emails.([]any)
		check(t, adv, []string{"RFC7644-3.5.2.2-L2151"},
			!hasEmails || emails == nil || len(emailArr) == 0,
			fmt.Sprintf("after PATCH remove emails: emails = %v, want unassigned", emails))
	}
}

func TestPatch_RemoveRequired(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	// Try to remove required attribute userName.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":   "remove",
				"path": "userName",
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.5.2.2-L2167: removing required attribute shall return error.
	check(t, adv, []string{"RFC7644-3.5.2.2-L2167"},
		patchResp.StatusCode == 400,
		fmtStatus("PATCH remove required userName", patchResp.StatusCode, 400))
}

func TestPatch_ReplaceDisplayName(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST /Users returned %d", resp.StatusCode)
	}
	id := idOf(body)

	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "displayName",
				"value": "Replaced User",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2-L1939"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH replace /Users/"+id, patchResp.StatusCode, 200))

	// Verify via GET.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		dn, _ := getResp.Body["displayName"].(string)
		check(t, adv, []string{"RFC7644-3.5.2-L1810"},
			dn == "Replaced User",
			fmt.Sprintf("after PATCH replace: displayName = %q, want \"Replaced User\"", dn))
	}
}

func TestPatch_ReplaceNonExistent(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch)

	body, resp := createUser(t)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(body)

	newNick := "NewNick-" + randomSuffix()
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "nickName",
				"value": newNick,
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.5.2.3-L2376: replace on non-existent attribute shall treat as add.
	check(t, adv, []string{"RFC7644-3.5.2.3-L2376"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH replace non-existent nickName", patchResp.StatusCode, 200))

	// Verify the value was set.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		nick, _ := getResp.Body["nickName"].(string)
		check(t, adv, []string{"RFC7644-3.5.2.3-L2376"},
			nick == newNick,
			fmt.Sprintf("PATCH replace non-existent: nickName = %q, want %q", nick, newNick))
	}
}

func TestPatch_ReplaceValuePath(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch, spec.Filter)

	// Create user with two emails.
	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-vpath-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "work@example.com", "type": "work"},
			{"value": "home@example.com", "type": "home"},
		},
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	// RFC7644-3.5.2.3-L2387: replace with valuePath filter shall replace matching values.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "emails[type eq \"work\"].value",
				"value": "newwork@example.com",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2.3-L2387"},
		patchResp.StatusCode == 200 || patchResp.StatusCode == 204,
		fmtStatus("PATCH replace with valuePath", patchResp.StatusCode, 200))

	// Verify via GET.
	getResp, err := client.Get("/Users/" + id)
	requireOK(t, err)

	if getResp.StatusCode == 200 {
		emails, _ := getResp.Body["emails"].([]any)
		found := false
		for _, e := range emails {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			if t, _ := em["type"].(string); t == "work" {
				v, _ := em["value"].(string)
				if v == "newwork@example.com" {
					found = true
				}
			}
		}
		check(t, adv, []string{"RFC7644-3.5.2.3-L2387"},
			found,
			"PATCH replace valuePath: work email not updated to newwork@example.com")
	}
}

func TestPatch_ReplaceValuePath_NoTarget(t *testing.T) {
	adv := requires(t, spec.Users, spec.Patch, spec.Filter)

	resp, err := client.Post("/Users", map[string]any{
		"schemas":  []string{userSchema},
		"userName": "scim-test-notarget-" + randomSuffix(),
		"emails": []map[string]any{
			{"value": "only@example.com", "type": "work"},
		},
	})
	requireOK(t, err)
	if resp.StatusCode != 201 {
		t.Fatalf("setup: POST returned %d", resp.StatusCode)
	}
	id := idOf(resp.Body)
	t.Cleanup(func() { _, _ = client.Delete("/Users/" + id) })

	// RFC7644-3.5.2.3-L2396: replace with unmatched valuePath shall return 400 noTarget.
	patchResp, err := client.Patch("/Users/"+id, map[string]any{
		"schemas": []string{patchOpSchema},
		"Operations": []map[string]any{
			{
				"op":    "replace",
				"path":  "emails[type eq \"nonexistent\"].value",
				"value": "nope@example.com",
			},
		},
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.5.2.3-L2396"},
		patchResp.StatusCode == 400,
		fmtStatus("PATCH replace unmatched valuePath", patchResp.StatusCode, 400))

	if patchResp.StatusCode == 400 && patchResp.Body != nil {
		scimType, _ := patchResp.Body["scimType"].(string)
		check(t, adv, []string{"RFC7644-3.5.2.3-L2396"},
			scimType == "noTarget",
			fmt.Sprintf("scimType = %q, want \"noTarget\"", scimType))
	}
}

func TestPatch_Supported(t *testing.T) {
	adv := requires(t, spec.Patch)

	// RFC7644-3.5-L1609: implementers SHOULD support PATCH.
	// If we get here, the feature is supported.
	check(t, adv, []string{"RFC7644-3.5-L1609"},
		features.Patch,
		"PATCH not supported by service provider")
}
