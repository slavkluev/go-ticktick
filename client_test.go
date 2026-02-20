package ticktick

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := NewClient("test-token")
	client.BaseURL = server.URL
	return client, server
}

func TestNewClient(t *testing.T) {
	client := NewClient("my-token")

	if client.AccessToken != "my-token" {
		t.Errorf("expected access token my-token, got %s", client.AccessToken)
	}
	if client.BaseURL != DefaultBaseURL {
		t.Errorf("expected base URL %s, got %s", DefaultBaseURL, client.BaseURL)
	}
	if client.HTTPClient != http.DefaultClient {
		t.Error("expected http.DefaultClient")
	}
}

func TestAuthorizationHeader(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	defer server.Close()

	client.get(context.Background(), "/test", &struct{}{})
}

func TestErrorResponse(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	})
	defer server.Close()

	err := client.get(context.Background(), "/test", &struct{}{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
	if apiErr.Body != "unauthorized" {
		t.Errorf("expected body unauthorized, got %s", apiErr.Body)
	}
}

func TestErrorString(t *testing.T) {
	err := &Error{StatusCode: 404, Body: "not found"}
	expected := "ticktick: HTTP 404: not found"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
