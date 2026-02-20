package ticktick

import (
	"context"
	"fmt"
	"net/url"
)

// GetProjects returns all projects for the authenticated user.
func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	var projects []Project
	if err := c.get(ctx, "/open/v1/project", &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject retrieves a project by ID.
func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	path := fmt.Sprintf("/open/v1/project/%s", url.PathEscape(projectID))
	var project Project
	if err := c.get(ctx, path, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// GetProjectData retrieves a project along with its tasks and columns.
func (c *Client) GetProjectData(ctx context.Context, projectID string) (*ProjectData, error) {
	path := fmt.Sprintf("/open/v1/project/%s/data", url.PathEscape(projectID))
	var data ProjectData
	if err := c.get(ctx, path, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	var project Project
	if err := c.post(ctx, "/open/v1/project", req, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(ctx context.Context, projectID string, req *UpdateProjectRequest) (*Project, error) {
	path := fmt.Sprintf("/open/v1/project/%s", url.PathEscape(projectID))
	var project Project
	if err := c.post(ctx, path, req, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteProject deletes a project.
func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	path := fmt.Sprintf("/open/v1/project/%s", url.PathEscape(projectID))
	return c.delete(ctx, path)
}
