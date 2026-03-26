package scim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ContentType = "application/scim+json"

// Auth applies authentication to an HTTP request.
type Auth interface {
	Apply(req *http.Request)
}

// BasicAuth authenticates with username and password.
type BasicAuth struct {
	User string
	Pass string
}

func (a BasicAuth) Apply(req *http.Request) {
	req.SetBasicAuth(a.User, a.Pass)
}

// BearerAuth authenticates with a bearer token.
type BearerAuth struct {
	Token string
}

func (a BearerAuth) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
}

// Client is a thin HTTP client for SCIM endpoints.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Auth       Auth
}

func (c *Client) Delete(path string) (*Response, error) {
	return c.Do(http.MethodDelete, path, nil)
}

// Do executes an HTTP request against the SCIM server.
func (c *Client) Do(method, path string, body map[string]any) (*Response, error) {
	return c.DoWithHeaders(method, path, body, nil)
}

// DoWithHeaders executes an HTTP request with additional headers.
func (c *Client) DoWithHeaders(method, path string, body map[string]any, headers map[string]string) (*Response, error) {
	url := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", ContentType)
	if body != nil {
		req.Header.Set("Content-Type", ContentType)
	}
	if c.Auth != nil {
		c.Auth.Apply(req)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	r := &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		RawBody:    raw,
	}

	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &r.Body); err != nil {
			// Not all responses are JSON (e.g. 204 No Content).
			// Store raw body but leave Body nil.
			r.Body = nil
		}
	}

	return r, nil
}

func (c *Client) Get(path string) (*Response, error) {
	return c.Do(http.MethodGet, path, nil)
}

func (c *Client) Patch(path string, body map[string]any) (*Response, error) {
	return c.Do(http.MethodPatch, path, body)
}

func (c *Client) Post(path string, body map[string]any) (*Response, error) {
	return c.Do(http.MethodPost, path, body)
}

func (c *Client) Put(path string, body map[string]any) (*Response, error) {
	return c.Do(http.MethodPut, path, body)
}

// Response wraps an HTTP response with a parsed JSON body.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       map[string]any
	RawBody    []byte
}
