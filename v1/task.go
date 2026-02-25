package ticktick

import (
	"context"
	"fmt"
	"net/url"
)

// GetTask retrieves a task by project ID and task ID.
func (c *Client) GetTask(ctx context.Context, projectID, taskID string) (*Task, error) {
	path := fmt.Sprintf("/open/v1/project/%s/task/%s", url.PathEscape(projectID), url.PathEscape(taskID))

	var task Task

	if err := c.get(ctx, path, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
	var task Task

	if err := c.post(ctx, "/open/v1/task", req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask updates an existing task.
func (c *Client) UpdateTask(ctx context.Context, taskID string, req *UpdateTaskRequest) (*Task, error) {
	path := fmt.Sprintf("/open/v1/task/%s", url.PathEscape(taskID))

	var task Task

	if err := c.post(ctx, path, req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CompleteTask marks a task as complete.
func (c *Client) CompleteTask(ctx context.Context, projectID, taskID string) error {
	path := fmt.Sprintf("/open/v1/project/%s/task/%s/complete", url.PathEscape(projectID), url.PathEscape(taskID))

	return c.post(ctx, path, nil, nil)
}

// DeleteTask deletes a task.
func (c *Client) DeleteTask(ctx context.Context, projectID, taskID string) error {
	path := fmt.Sprintf("/open/v1/project/%s/task/%s", url.PathEscape(projectID), url.PathEscape(taskID))

	return c.delete(ctx, path)
}
