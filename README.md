# go-ticktick

Go client library for the [TickTick Open API](https://developer.ticktick.com/manage).

## Installation

```bash
go get github.com/slavkluev/go-ticktick
```

## Authentication

The library requires an OAuth2 access token. To obtain one:

1. Register your application at the [TickTick Developer Center](https://developer.ticktick.com/manage) to get a `client_id` and `client_secret`.

2. Redirect the user to the authorization page:
   ```
   https://ticktick.com/oauth/authorize?scope=tasks:write%20tasks:read&client_id=YOUR_CLIENT_ID&state=STATE&redirect_uri=YOUR_REDIRECT_URI&response_type=code
   ```

3. After the user grants access, TickTick redirects back to your `redirect_uri` with a `code` parameter.

4. Exchange the code for an access token by making a POST request to `https://ticktick.com/oauth/token` with Basic Auth (`client_id` as username, `client_secret` as password) and form body:
   ```
   code=AUTHORIZATION_CODE&grant_type=authorization_code&scope=tasks:write%20tasks:read&redirect_uri=YOUR_REDIRECT_URI
   ```

   The response contains the `access_token`.

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
		Priority:  ticktick.Int(3),
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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
