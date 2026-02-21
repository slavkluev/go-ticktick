package ticktick_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/slavkluev/go-ticktick"
)

func TestGetProjects(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode([]ticktick.Project{
			{ID: "proj1", Name: "Project 1"},
			{ID: "proj2", Name: "Project 2"},
		})
	})
	defer server.Close()

	projects, err := client.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	if projects[0].ID != "proj1" {
		t.Errorf("expected project ID proj1, got %s", projects[0].ID)
	}
}

func TestGetProject(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj1", Name: "Project 1"})
	})
	defer server.Close()

	project, err := client.GetProject(context.Background(), "proj1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.ID != "proj1" {
		t.Errorf("expected project ID proj1, got %s", project.ID)
	}

	if project.Name != "Project 1" {
		t.Errorf("expected name Project 1, got %s", project.Name)
	}
}

func TestGetProjectData(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1/data" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(ticktick.ProjectData{
			Project: ticktick.Project{ID: "proj1", Name: "Project 1"},
			Tasks:   []ticktick.Task{{ID: "task1", Title: "Task 1"}},
			Columns: []ticktick.Column{{ID: "col1", Name: "Column 1"}},
		})
	})
	defer server.Close()

	data, err := client.GetProjectData(context.Background(), "proj1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Project.ID != "proj1" {
		t.Errorf("expected project ID proj1, got %s", data.Project.ID)
	}

	if len(data.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(data.Tasks))
	}

	if len(data.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(data.Columns))
	}
}

func TestCreateProject(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ticktick.CreateProjectRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.Name != "New Project" {
			t.Errorf("expected name New Project, got %s", req.Name)
		}

		if req.Color == nil || *req.Color != "#F18181" {
			t.Errorf("expected color #F18181, got %v", req.Color)
		}

		if req.SortOrder == nil || *req.SortOrder != 100 {
			t.Errorf("expected sortOrder 100, got %v", req.SortOrder)
		}

		if req.ViewMode == nil || *req.ViewMode != "list" {
			t.Errorf("expected viewMode list, got %v", req.ViewMode)
		}

		if req.Kind == nil || *req.Kind != "TASK" {
			t.Errorf("expected kind TASK, got %v", req.Kind)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj-new", Name: "New Project"})
	})
	defer server.Close()

	project, err := client.CreateProject(context.Background(), &ticktick.CreateProjectRequest{
		Name:      "New Project",
		Color:     ticktick.String("#F18181"),
		SortOrder: ticktick.Int64(100),
		ViewMode:  ticktick.String(ticktick.ViewModeList),
		Kind:      ticktick.String(ticktick.ProjectKindTask),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.ID != "proj-new" {
		t.Errorf("expected project ID proj-new, got %s", project.ID)
	}
}

func TestUpdateProject(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req ticktick.UpdateProjectRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.Name == nil || *req.Name != "Updated" {
			t.Errorf("expected name Updated, got %v", req.Name)
		}

		if req.Color == nil || *req.Color != "#00FF00" {
			t.Errorf("expected color #00FF00, got %v", req.Color)
		}

		if req.SortOrder == nil || *req.SortOrder != 200 {
			t.Errorf("expected sortOrder 200, got %v", req.SortOrder)
		}

		if req.ViewMode == nil || *req.ViewMode != "kanban" {
			t.Errorf("expected viewMode kanban, got %v", req.ViewMode)
		}

		if req.Kind == nil || *req.Kind != "NOTE" {
			t.Errorf("expected kind NOTE, got %v", req.Kind)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj1", Name: "Updated"})
	})
	defer server.Close()

	project, err := client.UpdateProject(context.Background(), "proj1", &ticktick.UpdateProjectRequest{
		Name:      ticktick.String("Updated"),
		Color:     ticktick.String("#00FF00"),
		SortOrder: ticktick.Int64(200),
		ViewMode:  ticktick.String(ticktick.ViewModeKanban),
		Kind:      ticktick.String(ticktick.ProjectKindNote),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", project.Name)
	}
}

func TestDeleteProject(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/open/v1/project/proj1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.DeleteProject(context.Background(), "proj1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateProjectOmitEmpty(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any

		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if raw["name"] != "Minimal" {
			t.Errorf("expected name Minimal, got %v", raw["name"])
		}

		for _, key := range []string{"color", "sortOrder", "viewMode", "kind"} {
			if _, ok := raw[key]; ok {
				t.Errorf("expected key %q to be omitted, but it was present", key)
			}
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "p1", Name: "Minimal"})
	})
	defer server.Close()

	_, err := client.CreateProject(context.Background(), &ticktick.CreateProjectRequest{
		Name: "Minimal",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateProjectOmitEmpty(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any

		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		for _, key := range []string{"name", "color", "sortOrder", "viewMode", "kind"} {
			if _, ok := raw[key]; ok {
				t.Errorf("expected key %q to be omitted, but it was present", key)
			}
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj1", Name: "Unchanged"})
	})
	defer server.Close()

	_, err := client.UpdateProject(context.Background(), "proj1", &ticktick.UpdateProjectRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetProjectPathEscaping(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		expected := "/open/v1/project/proj%2F1"
		if r.URL.RawPath != expected {
			t.Errorf("expected raw path %s, got %s", expected, r.URL.RawPath)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj/1", Name: "Project"})
	})
	defer server.Close()

	project, err := client.GetProject(context.Background(), "proj/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.ID != "proj/1" {
		t.Errorf("expected project ID proj/1, got %s", project.ID)
	}
}

func TestGetProjectError(t *testing.T) {
	client, server := setupTestClient(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("project not found"))
	})
	defer server.Close()

	_, err := client.GetProject(context.Background(), "proj1")
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

	if apiErr.Body != "project not found" {
		t.Errorf("expected body \"project not found\", got %s", apiErr.Body)
	}
}
