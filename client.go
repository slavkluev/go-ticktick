package ticktick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// DefaultBaseURL is the default base URL of the TickTick API.
const DefaultBaseURL = "https://api.ticktick.com"

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

func (c *Client) do(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req) //nolint:gosec // G704: URL is constructed from client-configured baseURL
	if err != nil {
		return err
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)

		return &Error{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	if v != nil {
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if errors.Is(decErr, io.EOF) {
			return &Error{
				StatusCode: resp.StatusCode,
				Body:       "empty response body",
			}
		}

		if decErr != nil {
			return decErr
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	return c.do(req, v)
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

	return c.do(req, v)
}

func (c *Client) delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
