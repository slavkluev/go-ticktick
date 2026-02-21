# go-ticktick

[![CI](https://github.com/slavkluev/go-ticktick/actions/workflows/ci.yml/badge.svg)](https://github.com/slavkluev/go-ticktick/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/slavkluev/go-ticktick/graph/badge.svg)](https://codecov.io/gh/slavkluev/go-ticktick)
[![Go Report Card](https://goreportcard.com/badge/github.com/slavkluev/go-ticktick)](https://goreportcard.com/report/github.com/slavkluev/go-ticktick)
[![Go Reference](https://pkg.go.dev/badge/github.com/slavkluev/go-ticktick.svg)](https://pkg.go.dev/github.com/slavkluev/go-ticktick)
[![Release](https://img.shields.io/github/v/release/slavkluev/go-ticktick.svg)](https://github.com/slavkluev/go-ticktick/releases/)
[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/slavkluev/go-ticktick/blob/main/LICENSE)

Go client library for the [TickTick Open API](https://developer.ticktick.com/manage).

## Installation

```bash
go get github.com/slavkluev/go-ticktick
```

## Authentication

The library requires an OAuth2 access token. TickTick uses the [Authorization Code](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1) flow.

### 1. Register your application

1. Open the [TickTick Developer Center](https://developer.ticktick.com/manage) and sign in.
2. Click **New App**, fill in **Name** (e.g. "My App") and click **Add**.
3. Click **Edit** on your app to open **App Setting**.
4. Copy your **Client ID** and **Client Secret**.
5. Set the **OAuth redirect URL** to a URL your app will listen on (e.g. `http://localhost:8080/callback`) and click **Save**.

### 2. Redirect user to authorization page

Build the authorization URL and redirect the user:

```
https://ticktick.com/oauth/authorize?client_id=YOUR_CLIENT_ID&scope=tasks:write%20tasks:read&redirect_uri=YOUR_REDIRECT_URI&response_type=code
```

| Parameter       | Description                                                                                  |
|-----------------|----------------------------------------------------------------------------------------------|
| `client_id`     | Your application ID                                                                          |
| `scope`         | `tasks:read` for read access, `tasks:write` for write access (both cover tasks and projects) |
| `state`         | Any value, passed back to redirect URL as-is (optional)                                      |
| `redirect_uri`  | Your registered redirect URL                                                                 |
| `response_type` | Must be `code`                                                                               |

### 3. Receive the authorization code

After the user grants access, TickTick redirects to your `redirect_uri` with query parameters:

| Parameter | Description                                      |
|-----------|--------------------------------------------------|
| `code`    | Authorization code to exchange for a token       |
| `state`   | The same `state` value from step 2 (if provided) |

### 4. Exchange code for access token

Exchange the authorization code for an access token:

```bash
curl -X POST https://ticktick.com/oauth/token \
  -u "YOUR_CLIENT_ID:YOUR_CLIENT_SECRET" \
  -d "code=AUTHORIZATION_CODE&grant_type=authorization_code&scope=tasks:write%20tasks:read&redirect_uri=YOUR_REDIRECT_URI"
```

The response JSON contains the `access_token`:

```json
{
  "access_token": "your-access-token",
  "token_type": "bearer",
  "expires_in": 15551999,
  "scope": "tasks:read tasks:write"
}
```

### 5. Use the token

Pass the token when creating the client:

```go
client := ticktick.NewClient("your-access-token")
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/slavkluev/go-ticktick"
)

func main() {
	client := ticktick.NewClient("your-access-token")
	ctx := context.Background()

	// List all projects
	projects, err := client.GetProjects(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range projects {
		fmt.Printf("Project: %s (%s)\n", p.Name, p.ID)
	}

	// Create a task
	task, err := client.CreateTask(ctx, &ticktick.CreateTaskRequest{
		Title:     "Buy groceries",
		ProjectID: projects[0].ID,
		Priority:  ticktick.Int(ticktick.PriorityMedium),
		IsAllDay:  ticktick.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created task: %s\n", task.ID)

	// Complete the task
	err = client.CompleteTask(ctx, task.ProjectID, task.ID)
	if err != nil {
		log.Fatal(err)
	}
}
```

## API

### Client

```go
client := ticktick.NewClient("access-token")

// Customize HTTP client or base URL
client = ticktick.NewClient("access-token",
	ticktick.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	ticktick.WithBaseURL("https://custom-api-url.example.com"),
)
```

### Tasks

| Method                                        | Description                       |
|-----------------------------------------------|-----------------------------------|
| `GetTask(ctx, projectID, taskID)`             | Get a task by project and task ID |
| `CreateTask(ctx, *CreateTaskRequest)`         | Create a new task                 |
| `UpdateTask(ctx, taskID, *UpdateTaskRequest)` | Update an existing task           |
| `CompleteTask(ctx, projectID, taskID)`        | Mark a task as complete           |
| `DeleteTask(ctx, projectID, taskID)`          | Delete a task                     |

### Projects

| Method                                                 | Description                              |
|--------------------------------------------------------|------------------------------------------|
| `GetProjects(ctx)`                                     | List all projects                        |
| `GetProject(ctx, projectID)`                           | Get a project by ID                      |
| `GetProjectData(ctx, projectID)`                       | Get a project with its tasks and columns |
| `CreateProject(ctx, *CreateProjectRequest)`            | Create a new project                     |
| `UpdateProject(ctx, projectID, *UpdateProjectRequest)` | Update an existing project               |
| `DeleteProject(ctx, projectID)`                        | Delete a project                         |

### Error Handling

API errors are returned as `*ticktick.Error` with the HTTP status code and response body:

```go
import "errors"

task, err := client.GetTask(ctx, "proj1", "task1")
if err != nil {
	var apiErr *ticktick.Error
	if errors.As(err, &apiErr) {
		fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Body)
	}
}
```

### Constants

The library provides constants for common field values:

```go
// Task priority
ticktick.PriorityNone      // 0
ticktick.PriorityLow       // 1
ticktick.PriorityMedium    // 3
ticktick.PriorityHigh      // 5

// Task status
ticktick.TaskStatusNormal       // 0
ticktick.TaskStatusCompleted    // 2

// Checklist item status
ticktick.ChecklistStatusNormal       // 0
ticktick.ChecklistStatusCompleted    // 1

// Project view mode
ticktick.ViewModeList       // "list"
ticktick.ViewModeKanban     // "kanban"
ticktick.ViewModeTimeline   // "timeline"

// Project kind
ticktick.ProjectKindTask    // "TASK"
ticktick.ProjectKindNote    // "NOTE"

// Task kind
ticktick.TaskKindText       // "TEXT"
ticktick.TaskKindNote       // "NOTE"
ticktick.TaskKindChecklist  // "CHECKLIST"

// Project permission
ticktick.PermissionRead     // "read"
ticktick.PermissionWrite    // "write"
ticktick.PermissionComment  // "comment"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
