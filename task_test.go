package ticktick_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/slavkluev/go-ticktick"
)

func TestGetTask(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1/task/task1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "task1", ProjectID: "proj1", Title: "Test Task"})
	})
	defer server.Close()

	task, err := client.GetTask(context.Background(), "proj1", "task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "task1" {
		t.Errorf("expected task ID task1, got %s", task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("expected title Test Task, got %s", task.Title)
	}
}

func TestCreateTask(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/task" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var req ticktick.CreateTaskRequest

		json.NewDecoder(r.Body).Decode(&req)

		if req.Title != "New Task" {
			t.Errorf("expected title New Task, got %s", req.Title)
		}

		if req.ProjectID != "proj1" {
			t.Errorf("expected projectId proj1, got %s", req.ProjectID)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "new1", ProjectID: "proj1", Title: "New Task"})
	})
	defer server.Close()

	task, err := client.CreateTask(context.Background(), &ticktick.CreateTaskRequest{
		Title:     "New Task",
		ProjectID: "proj1",
		Priority:  ticktick.Int(3),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "new1" {
		t.Errorf("expected task ID new1, got %s", task.ID)
	}
}

func TestUpdateTask(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/task/task1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ticktick.UpdateTaskRequest

		json.NewDecoder(r.Body).Decode(&req)

		if req.ID != "task1" {
			t.Errorf("expected id task1, got %s", req.ID)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "task1", ProjectID: "proj1", Title: "Updated"})
	})
	defer server.Close()

	task, err := client.UpdateTask(context.Background(), "task1", &ticktick.UpdateTaskRequest{
		ID:        "task1",
		ProjectID: "proj1",
		Title:     ticktick.String("Updated"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.Title != "Updated" {
		t.Errorf("expected title Updated, got %s", task.Title)
	}
}

func TestCompleteTask(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1/task/task1/complete" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.CompleteTask(context.Background(), "proj1", "task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1/task/task1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.DeleteTask(context.Background(), "proj1", "task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
