package ticktick_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

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
	startDate := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	dueDate := time.Date(2024, 1, 20, 18, 0, 0, 0, time.UTC)

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

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.Title != "New Task" {
			t.Errorf("expected title New Task, got %s", req.Title)
		}

		if req.ProjectID != "proj1" {
			t.Errorf("expected projectId proj1, got %s", req.ProjectID)
		}

		if req.Content == nil || *req.Content != "task content" {
			t.Errorf("expected content task content, got %v", req.Content)
		}

		if req.Desc == nil || *req.Desc != "task description" {
			t.Errorf("expected desc task description, got %v", req.Desc)
		}

		if req.IsAllDay == nil || *req.IsAllDay != true {
			t.Errorf("expected isAllDay true, got %v", req.IsAllDay)
		}

		if req.StartDate == nil || !req.StartDate.Equal(startDate) {
			t.Errorf("expected startDate %v, got %v", startDate, req.StartDate)
		}

		if req.DueDate == nil || !req.DueDate.Equal(dueDate) {
			t.Errorf("expected dueDate %v, got %v", dueDate, req.DueDate)
		}

		if req.TimeZone == nil || *req.TimeZone != "America/New_York" {
			t.Errorf("expected timeZone America/New_York, got %v", req.TimeZone)
		}

		if len(req.Reminders) != 2 || req.Reminders[0] != "TRIGGER:P0DT9H0M0S" || req.Reminders[1] != "TRIGGER:PT0S" {
			t.Errorf("unexpected reminders: %v", req.Reminders)
		}

		if req.RepeatFlag == nil || *req.RepeatFlag != "RRULE:FREQ=DAILY;INTERVAL=1" {
			t.Errorf("expected repeatFlag RRULE:FREQ=DAILY;INTERVAL=1, got %v", req.RepeatFlag)
		}

		if req.Priority == nil || *req.Priority != ticktick.PriorityHigh {
			t.Errorf("expected priority %d, got %v", ticktick.PriorityHigh, req.Priority)
		}

		if req.SortOrder == nil || *req.SortOrder != 12345 {
			t.Errorf("expected sortOrder 12345, got %v", req.SortOrder)
		}

		if len(req.Items) != 1 || req.Items[0].Title != "Subtask 1" {
			t.Errorf("unexpected items: %v", req.Items)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "new1", ProjectID: "proj1", Title: "New Task"})
	})
	defer server.Close()

	task, err := client.CreateTask(context.Background(), &ticktick.CreateTaskRequest{
		Title:      "New Task",
		ProjectID:  "proj1",
		Content:    ticktick.String("task content"),
		Desc:       ticktick.String("task description"),
		IsAllDay:   ticktick.Bool(true),
		StartDate:  ticktick.NewTime(startDate),
		DueDate:    ticktick.NewTime(dueDate),
		TimeZone:   ticktick.String("America/New_York"),
		Reminders:  []string{"TRIGGER:P0DT9H0M0S", "TRIGGER:PT0S"},
		RepeatFlag: ticktick.String("RRULE:FREQ=DAILY;INTERVAL=1"),
		Priority:   ticktick.Int(ticktick.PriorityHigh),
		SortOrder:  ticktick.Int64(12345),
		Items:      []ticktick.CreateChecklistItemRequest{{Title: "Subtask 1"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "new1" {
		t.Errorf("expected task ID new1, got %s", task.ID)
	}
}

func TestUpdateTask(t *testing.T) {
	startDate := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
	dueDate := time.Date(2024, 2, 5, 17, 0, 0, 0, time.UTC)

	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/task/task1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ticktick.UpdateTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.ID != "task1" {
			t.Errorf("expected id task1, got %s", req.ID)
		}

		if req.ProjectID != "proj1" {
			t.Errorf("expected projectId proj1, got %s", req.ProjectID)
		}

		if req.Title == nil || *req.Title != "Updated" {
			t.Errorf("expected title Updated, got %v", req.Title)
		}

		if req.Content == nil || *req.Content != "updated content" {
			t.Errorf("expected content updated content, got %v", req.Content)
		}

		if req.Desc == nil || *req.Desc != "updated desc" {
			t.Errorf("expected desc updated desc, got %v", req.Desc)
		}

		if req.IsAllDay == nil || *req.IsAllDay != false {
			t.Errorf("expected isAllDay false, got %v", req.IsAllDay)
		}

		if req.StartDate == nil || !req.StartDate.Equal(startDate) {
			t.Errorf("expected startDate %v, got %v", startDate, req.StartDate)
		}

		if req.DueDate == nil || !req.DueDate.Equal(dueDate) {
			t.Errorf("expected dueDate %v, got %v", dueDate, req.DueDate)
		}

		if req.TimeZone == nil || *req.TimeZone != "Europe/London" {
			t.Errorf("expected timeZone Europe/London, got %v", req.TimeZone)
		}

		if len(req.Reminders) != 1 || req.Reminders[0] != "TRIGGER:PT0S" {
			t.Errorf("unexpected reminders: %v", req.Reminders)
		}

		if req.RepeatFlag == nil || *req.RepeatFlag != "RRULE:FREQ=WEEKLY;INTERVAL=2" {
			t.Errorf("expected repeatFlag RRULE:FREQ=WEEKLY;INTERVAL=2, got %v", req.RepeatFlag)
		}

		if req.Priority == nil || *req.Priority != ticktick.PriorityMedium {
			t.Errorf("expected priority %d, got %v", ticktick.PriorityMedium, req.Priority)
		}

		if req.SortOrder == nil || *req.SortOrder != 99999 {
			t.Errorf("expected sortOrder 99999, got %v", req.SortOrder)
		}

		if len(req.Items) != 1 || req.Items[0].Title != "Updated subtask" {
			t.Errorf("unexpected items: %v", req.Items)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "task1", ProjectID: "proj1", Title: "Updated"})
	})
	defer server.Close()

	task, err := client.UpdateTask(context.Background(), "task1", &ticktick.UpdateTaskRequest{
		ID:         "task1",
		ProjectID:  "proj1",
		Title:      ticktick.String("Updated"),
		Content:    ticktick.String("updated content"),
		Desc:       ticktick.String("updated desc"),
		IsAllDay:   ticktick.Bool(false),
		StartDate:  ticktick.NewTime(startDate),
		DueDate:    ticktick.NewTime(dueDate),
		TimeZone:   ticktick.String("Europe/London"),
		Reminders:  []string{"TRIGGER:PT0S"},
		RepeatFlag: ticktick.String("RRULE:FREQ=WEEKLY;INTERVAL=2"),
		Priority:   ticktick.Int(ticktick.PriorityMedium),
		SortOrder:  ticktick.Int64(99999),
		Items:      []ticktick.CreateChecklistItemRequest{{Title: "Updated subtask"}},
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

func TestCreateTaskOmitEmpty(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any

		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if raw["title"] != "Minimal" {
			t.Errorf("expected title Minimal, got %v", raw["title"])
		}

		if raw["projectId"] != "proj1" {
			t.Errorf("expected projectId proj1, got %v", raw["projectId"])
		}

		for _, key := range []string{
			"content", "desc", "isAllDay", "startDate", "dueDate",
			"timeZone", "reminders", "repeatFlag", "priority", "sortOrder", "items",
		} {
			if _, ok := raw[key]; ok {
				t.Errorf("expected key %q to be omitted, but it was present", key)
			}
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "t1", Title: "Minimal"})
	})
	defer server.Close()

	_, err := client.CreateTask(context.Background(), &ticktick.CreateTaskRequest{
		Title:     "Minimal",
		ProjectID: "proj1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateTaskOmitEmpty(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any

		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if raw["id"] != "task1" {
			t.Errorf("expected id task1, got %v", raw["id"])
		}

		if raw["projectId"] != "proj1" {
			t.Errorf("expected projectId proj1, got %v", raw["projectId"])
		}

		for _, key := range []string{
			"title", "content", "desc", "isAllDay", "startDate", "dueDate",
			"timeZone", "reminders", "repeatFlag", "priority", "sortOrder", "items",
		} {
			if _, ok := raw[key]; ok {
				t.Errorf("expected key %q to be omitted, but it was present", key)
			}
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "task1"})
	})
	defer server.Close()

	_, err := client.UpdateTask(context.Background(), "task1", &ticktick.UpdateTaskRequest{
		ID:        "task1",
		ProjectID: "proj1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTaskPathEscaping(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expected := "/open/v1/project/proj%2F1/task/task%3F2"
		if r.URL.RawPath != expected {
			t.Errorf("expected raw path %s, got %s", expected, r.URL.RawPath)
		}

		json.NewEncoder(w).Encode(ticktick.Task{ID: "task?2", ProjectID: "proj/1"})
	})
	defer server.Close()

	task, err := client.GetTask(context.Background(), "proj/1", "task?2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "task?2" {
		t.Errorf("expected task ID task?2, got %s", task.ID)
	}
}

func TestGetTaskError(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("task not found"))
	})
	defer server.Close()

	_, err := client.GetTask(context.Background(), "proj1", "task1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *ticktick.Error

	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ticktick.Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}

	if apiErr.Body != "task not found" {
		t.Errorf("expected body \"task not found\", got %s", apiErr.Body)
	}
}
