package ticktick_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slavkluev/go-ticktick"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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
	const markerHeader = "X-Custom-Client"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(markerHeader) != "true" {
			t.Error("expected custom client marker header")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	customClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set(markerHeader, "true")

			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	client := ticktick.NewClient("token",
		ticktick.WithBaseURL(server.URL),
		ticktick.WithHTTPClient(customClient),
	)

	_, err := client.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

	_, err := client.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

func TestErrorStatusBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		wantError bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"299 upper success bound", 299, false},
		{"300 redirect", 300, true},
		{"400 bad request", http.StatusBadRequest, true},
		{"429 rate limit", http.StatusTooManyRequests, true},
		{"500 server error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte("[]"))
			})
			defer server.Close()

			_, err := client.GetProjects(context.Background())

			if !tt.wantError {
				if err != nil {
					t.Fatalf("unexpected error for status %d: %v", tt.status, err)
				}

				return
			}

			if err == nil {
				t.Fatalf("expected error for status %d, got nil", tt.status)
			}

			var apiErr *ticktick.Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *ticktick.Error, got %T", err)
			}

			if apiErr.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, apiErr.StatusCode)
			}
		})
	}
}
