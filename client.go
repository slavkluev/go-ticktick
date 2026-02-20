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
	// HTTPClient is the HTTP client used for API requests.
	HTTPClient *http.Client

	// BaseURL is the base URL of the TickTick API.
	BaseURL string

	// AccessToken is the OAuth2 bearer token.
	AccessToken string
}

// NewClient creates a new TickTick API client with the given access token.
func NewClient(accessToken string) *Client {
	return &Client{
		HTTPClient:  http.DefaultClient,
		BaseURL:     DefaultBaseURL,
		AccessToken: accessToken,
	}
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
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, reqBody)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
