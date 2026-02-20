package ticktick_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slavkluev/go-ticktick"
)

func setupTestClient(handler http.HandlerFunc) (*ticktick.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := ticktick.NewClient("test-token", ticktick.WithBaseURL(server.URL))

	return client, server
}

func TestNewClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer my-token" {
			t.Errorf("expected Bearer my-token, got %s", auth)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := ticktick.NewClient("my-token", ticktick.WithBaseURL(server.URL))

	_, err := client.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	called := false
	customClient := &http.Client{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := ticktick.NewClient("token",
		ticktick.WithBaseURL(server.URL),
		ticktick.WithHTTPClient(customClient),
	)

	_, _ = client.GetProjects(context.Background())

	if !called {
		t.Error("expected custom HTTP client to be used")
	}
}

func TestAuthorizationHeader(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})
	defer server.Close()

	_, _ = client.GetProjects(context.Background())
}

func TestErrorResponse(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	})
	defer server.Close()

	_, err := client.GetProjects(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *ticktick.Error

	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ticktick.Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}

	if apiErr.Body != "unauthorized" {
		t.Errorf("expected body unauthorized, got %s", apiErr.Body)
	}
}

func TestErrorString(t *testing.T) {
	err := &ticktick.Error{StatusCode: 404, Body: "not found"}
	expected := "ticktick: HTTP 404: not found"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
