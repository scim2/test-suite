package compliance

import (
	"fmt"
	"testing"

	"github.com/scim2/test-suite/spec"
)

const bulkSchema = "urn:ietf:params:scim:api:messages:2.0:BulkRequest"

func TestBulk_CircularReferences(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users, spec.Groups)

	// RFC7644-3.7.1-L2814: must try to resolve circular cross-references.
	// Create two groups that reference each other via bulkId.
	resp, err := client.Post("/Bulk", map[string]any{
		"schemas": []string{bulkSchema},
		"Operations": []map[string]any{
			{
				"method": "POST",
				"path":   "/Groups",
				"bulkId": "groupA",
				"data": map[string]any{
					"schemas":     []string{groupSchema},
					"displayName": "scim-test-circA-" + randomSuffix(),
					"members": []map[string]any{
						{"value": "bulkId:groupB", "type": "Group"},
					},
				},
			},
			{
				"method": "POST",
				"path":   "/Groups",
				"bulkId": "groupB",
				"data": map[string]any{
					"schemas":     []string{groupSchema},
					"displayName": "scim-test-circB-" + randomSuffix(),
					"members": []map[string]any{
						{"value": "bulkId:groupA", "type": "Group"},
					},
				},
			},
		},
	})
	requireOK(t, err)

	// Server should attempt to resolve, returning 200 (even if individual ops fail).
	check(t, adv, []string{"RFC7644-3.7.1-L2814"},
		resp.StatusCode == 200,
		fmtStatus("POST /Bulk with circular refs", resp.StatusCode, 200))

	// Clean up.
	if resp.StatusCode == 200 && resp.Body != nil {
		ops, _ := resp.Body["Operations"].([]any)
		for _, op := range ops {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			if loc, _ := o["location"].(string); loc != "" {
				_, _ = client.Delete(loc)
			}
		}
	}
}

func TestBulk_ContinuesOnPartialFailure(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users)

	userName := "scim-test-bulkdup-" + randomSuffix()

	resp, err := client.Post("/Bulk", map[string]any{
		"schemas":      []string{bulkSchema},
		"failOnErrors": 10,
		"Operations": []map[string]any{
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "ok-1",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": userName,
				},
			},
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "dup-1",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": userName,
				},
			},
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "ok-2",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": "scim-test-bulkok-" + randomSuffix(),
				},
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.7-L2779: must continue and disregard partial failures.
	check(t, adv, []string{"RFC7644-3.7-L2779"},
		resp.StatusCode == 200,
		fmtStatus("POST /Bulk (partial failure)", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		ops, _ := resp.Body["Operations"].([]any)

		check(t, adv, []string{"RFC7644-3.7-L2779"},
			len(ops) == 3,
			fmt.Sprintf("bulk returned %d operations, want 3 (all processed)", len(ops)))

		// Clean up.
		for _, op := range ops {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			if loc, _ := o["location"].(string); loc != "" {
				_, _ = client.Delete(loc)
			}
		}
	}
}

func TestBulk_ExceedLimits(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users)

	// RFC7644-3.7.4-L3499: exceeding limits must return 413.
	// Discover maxOperations.
	spcResp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	bulkCfg, _ := spcResp.Body["bulk"].(map[string]any)
	maxOps, _ := bulkCfg["maxOperations"].(float64)
	if maxOps == 0 {
		maxOps = 10
	}

	// Build more operations than allowed.
	ops := make([]map[string]any, int(maxOps)+1)
	for i := range ops {
		ops[i] = map[string]any{
			"method": "POST",
			"path":   "/Users",
			"bulkId": fmt.Sprintf("excess-%d", i),
			"data": map[string]any{
				"schemas":  []string{userSchema},
				"userName": fmt.Sprintf("scim-test-exceed-%d-%s", i, randomSuffix()),
			},
		}
	}

	resp, err := client.Post("/Bulk", map[string]any{
		"schemas":    []string{bulkSchema},
		"Operations": ops,
	})
	requireOK(t, err)

	check(t, adv, []string{"RFC7644-3.7.4-L3499"},
		resp.StatusCode == 413,
		fmtStatus("POST /Bulk exceeding maxOperations", resp.StatusCode, 413))
}

func TestBulk_MaxLimits(t *testing.T) {
	adv := requires(t, spec.Bulk)

	// RFC7644-3.7.4-L3495: must define max operations and max payload size.
	spcResp, err := client.Get("/ServiceProviderConfig")
	requireOK(t, err)

	bulkCfg, _ := spcResp.Body["bulk"].(map[string]any)
	maxOps, hasMaxOps := bulkCfg["maxOperations"]
	maxPayload, hasMaxPayload := bulkCfg["maxPayloadSize"]

	check(t, adv, []string{"RFC7644-3.7.4-L3495"},
		hasMaxOps && hasMaxPayload,
		fmt.Sprintf("bulk config missing limits: maxOperations=%v, maxPayloadSize=%v", maxOps, maxPayload))
}

func TestBulk_ResponseIncludesAllOps(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users)

	resp, err := client.Post("/Bulk", map[string]any{
		"schemas": []string{bulkSchema},
		"Operations": []map[string]any{
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "r1",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": "scim-test-bulkall1-" + randomSuffix(),
				},
			},
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "r2",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": "scim-test-bulkall2-" + randomSuffix(),
				},
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.7.3-L3201: response must include result of all processed ops.
	check(t, adv, []string{"RFC7644-3.7.3-L3201"},
		resp.StatusCode == 200,
		fmtStatus("POST /Bulk (response ops)", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		ops, _ := resp.Body["Operations"].([]any)

		check(t, adv, []string{"RFC7644-3.7.3-L3201"},
			len(ops) == 2,
			fmt.Sprintf("bulk response has %d operations, want 2", len(ops)))

		// RFC7644-3.7.3-L3206: each op must include HTTP response code.
		for i, op := range ops {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			status, _ := o["status"].(string)
			check(t, adv, []string{"RFC7644-3.7.3-L3206"},
				status != "",
				fmt.Sprintf("bulk operation %d missing status", i))
		}

		// Clean up.
		for _, op := range ops {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			if loc, _ := o["location"].(string); loc != "" {
				_, _ = client.Delete(loc)
			}
		}
	} else {
		// Record checks so requirements aren't left untested.
		check(t, adv, []string{"RFC7644-3.7.3-L3206"},
			false,
			fmtStatus("POST /Bulk (status codes)", resp.StatusCode, 200))
	}
}

func TestBulk_Returns200(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users)

	resp, err := client.Post("/Bulk", map[string]any{
		"schemas": []string{bulkSchema},
		"Operations": []map[string]any{
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "test-1",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": "scim-test-bulk-" + randomSuffix(),
				},
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.7-L2775: successful bulk job must return 200 OK.
	check(t, adv, []string{"RFC7644-3.7-L2775"},
		resp.StatusCode == 200,
		fmtStatus("POST /Bulk", resp.StatusCode, 200))

	// Clean up created resources.
	if resp.StatusCode == 200 && resp.Body != nil {
		ops, _ := resp.Body["Operations"].([]any)
		for _, op := range ops {
			o, ok := op.(map[string]any)
			if !ok {
				continue
			}
			loc, _ := o["location"].(string)
			if loc != "" {
				_, _ = client.Delete(loc)
			}
		}
	}
}

func TestBulk_ReturnsBulkId(t *testing.T) {
	adv := requires(t, spec.Bulk, spec.Users)

	resp, err := client.Post("/Bulk", map[string]any{
		"schemas": []string{bulkSchema},
		"Operations": []map[string]any{
			{
				"method": "POST",
				"path":   "/Users",
				"bulkId": "myBulkId-42",
				"data": map[string]any{
					"schemas":  []string{userSchema},
					"userName": "scim-test-bulkid-" + randomSuffix(),
				},
			},
		},
	})
	requireOK(t, err)

	// RFC7644-3.7-L2789: must return same bulkId with newly created resource.
	check(t, adv, []string{"RFC7644-3.7-L2789"},
		resp.StatusCode == 200,
		fmtStatus("POST /Bulk (bulkId)", resp.StatusCode, 200))

	if resp.StatusCode == 200 {
		ops, _ := resp.Body["Operations"].([]any)
		if len(ops) > 0 {
			first, _ := ops[0].(map[string]any)
			bulkId, _ := first["bulkId"].(string)
			check(t, adv, []string{"RFC7644-3.7-L2789"},
				bulkId == "myBulkId-42",
				fmt.Sprintf("bulkId = %q, want \"myBulkId-42\"", bulkId))

			// Clean up.
			if loc, _ := first["location"].(string); loc != "" {
				_, _ = client.Delete(loc)
			}
		}
	}
}
