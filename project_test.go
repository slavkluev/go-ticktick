package ticktick_test

import (
	"context"
	"encoding/json"
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

		json.NewDecoder(r.Body).Decode(&req)

		if req.Name != "New Project" {
			t.Errorf("expected name New Project, got %s", req.Name)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj-new", Name: "New Project"})
	})
	defer server.Close()

	project, err := client.CreateProject(context.Background(), &ticktick.CreateProjectRequest{
		Name:     "New Project",
		Color:    ticktick.String("#F18181"),
		ViewMode: ticktick.String("list"),
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

		json.NewDecoder(r.Body).Decode(&req)

		if req.Name == nil || *req.Name != "Updated" {
			t.Errorf("expected name Updated, got %v", req.Name)
		}

		json.NewEncoder(w).Encode(ticktick.Project{ID: "proj1", Name: "Updated"})
	})
	defer server.Close()

	project, err := client.UpdateProject(context.Background(), "proj1", &ticktick.UpdateProjectRequest{
		Name: ticktick.String("Updated"),
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
