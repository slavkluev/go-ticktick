// Copyright (c) 2026 Viacheslav Kliuev
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

/*
Package ticktick provides a Go client for the TickTick Open API.

# Usage

	import "github.com/slavkluev/go-ticktick"

Construct a new client with an OAuth2 access token, then call methods
to access the TickTick API. For example:

	client := ticktick.NewClient("your-access-token")

	// List all projects.
	projects, err := client.GetProjects(ctx)

	// Create a high-priority task.
	task, err := client.CreateTask(ctx, &ticktick.CreateTaskRequest{
		Title:     "Buy groceries",
		ProjectID: projects[0].ID,
		Priority:  ticktick.Int(ticktick.PriorityHigh),
	})

	// Mark it complete.
	err = client.CompleteTask(ctx, task.ProjectID, task.ID)

The client supports customization through functional options:

	client := ticktick.NewClient("your-access-token",
		ticktick.WithHTTPClient(httpClient),
		ticktick.WithBaseURL("https://api.dida365.com"),
	)

All methods accept a [context.Context] as the first parameter for
cancellation and timeouts.

# Authentication

The TickTick API uses OAuth2 Bearer tokens. Obtain an access token through
the Authorization Code flow described in the TickTick Developer Center at
https://developer.ticktick.com/api, then pass it to [NewClient].

# Creating and Updating Resources

Request types use pointer fields for optional values, allowing the API to
distinguish between unset fields and zero values. Helper functions [String],
[Int], [Int64], [Bool], and [NewTime] create the required pointers:

	req := &ticktick.CreateTaskRequest{
		Title:     "Weekly report",
		ProjectID: "project-id",
		Content:   ticktick.String("Status update"),
		Priority:  ticktick.Int(ticktick.PriorityMedium),
		DueDate:   ticktick.NewTime(time.Now().Add(24 * time.Hour)),
	}

# Error Handling

API errors are returned as [*Error] with the HTTP status code and
response body. Use [errors.As] to inspect them:

	task, err := client.GetTask(ctx, projectID, taskID)
	if err != nil {
		var apiErr *ticktick.Error
		if errors.As(err, &apiErr) {
			log.Printf("HTTP %d: %s", apiErr.StatusCode, apiErr.Body)
		}
	}
*/
package ticktick
