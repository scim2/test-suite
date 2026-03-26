package scim

import "fmt"

// supported checks if a feature object has "supported": true.
func supported(body map[string]any, key string) bool {
	obj, ok := body[key].(map[string]any)
	if !ok {
		return false
	}
	v, ok := obj["supported"].(bool)
	return ok && v
}

// Discover queries the service provider's configuration and resource
// type endpoints to determine supported features.
func (c *Client) Discover() (*Features, error) {
	f := &Features{}

	if err := c.discoverConfig(f); err != nil {
		return nil, err
	}
	if err := c.discoverResourceTypes(f); err != nil {
		return nil, err
	}

	return f, nil
}

func (c *Client) discoverConfig(f *Features) error {
	resp, err := c.Get("/ServiceProviderConfig")
	if err != nil {
		return fmt.Errorf("GET /ServiceProviderConfig: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET /ServiceProviderConfig: status %d", resp.StatusCode)
	}
	if resp.Body == nil {
		return fmt.Errorf("GET /ServiceProviderConfig: empty body")
	}

	f.Patch = supported(resp.Body, "patch")
	f.Bulk = supported(resp.Body, "bulk")
	f.Filter = supported(resp.Body, "filter")
	f.ChangePassword = supported(resp.Body, "changePassword")
	f.Sort = supported(resp.Body, "sort")
	f.ETag = supported(resp.Body, "etag")

	if bulk, ok := resp.Body["bulk"].(map[string]any); ok {
		if v, ok := bulk["maxOperations"].(float64); ok {
			f.BulkMaxOps = int(v)
		}
		if v, ok := bulk["maxPayloadSize"].(float64); ok {
			f.BulkMaxPayload = int(v)
		}
	}

	if filter, ok := resp.Body["filter"].(map[string]any); ok {
		if v, ok := filter["maxResults"].(float64); ok {
			f.FilterMaxResults = int(v)
		}
	}

	return nil
}

func (c *Client) discoverResourceTypes(f *Features) error {
	resp, err := c.Get("/ResourceTypes")
	if err != nil {
		return fmt.Errorf("GET /ResourceTypes: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GET /ResourceTypes: status %d", resp.StatusCode)
	}
	if resp.Body == nil {
		return fmt.Errorf("GET /ResourceTypes: empty body")
	}

	// ResourceTypes returns a ListResponse with a Resources array.
	resources, ok := resp.Body["Resources"].([]any)
	if !ok {
		return nil
	}

	for _, r := range resources {
		rt, ok := r.(map[string]any)
		if !ok {
			continue
		}
		name, _ := rt["name"].(string)
		switch name {
		case "User":
			f.Users = true
		case "Group":
			f.Groups = true
		case "TestResource":
			f.TestResource = true
		}
	}

	return nil
}

// Features represents the capabilities of a SCIM service provider,
// discovered via GET /ServiceProviderConfig and GET /ResourceTypes.
type Features struct {
	Patch            bool
	Bulk             bool
	BulkMaxOps       int
	BulkMaxPayload   int
	Filter           bool
	FilterMaxResults int
	ChangePassword   bool
	Sort             bool
	ETag             bool
	Users            bool
	Groups           bool
	TestResource     bool
}
