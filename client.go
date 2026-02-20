package ticktick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// DefaultBaseURL is the default base URL of the TickTick API.
	DefaultBaseURL = "https://api.ticktick.com"

	// TimeLayout is the time format used by the TickTick API.
	TimeLayout = "2006-01-02T15:04:05-0700"
)

// Client manages communication with the TickTick Open API.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new TickTick API client with the given access token.
func NewClient(accessToken string, opts ...Option) *Client {
	c := &Client{
		httpClient:  http.DefaultClient,
		baseURL:     DefaultBaseURL,
		accessToken: accessToken,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Error represents an error response from the TickTick API.
type Error struct {
	StatusCode int
	Body       string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ticktick: HTTP %d: %s", e.StatusCode, e.Body)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req) //nolint:gosec // G704: URL is constructed from client-configured baseURL
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		return nil, &Error{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	return resp, nil
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) post(ctx context.Context, path string, body any, v any) error {
	var reqBody io.Reader

	if body != nil {
		var buf bytes.Buffer

		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}

		reqBody = &buf
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

func (c *Client) delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
