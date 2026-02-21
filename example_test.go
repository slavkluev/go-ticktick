package ticktick_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/slavkluev/go-ticktick"
)

func ExampleNewClient() {
	client := ticktick.NewClient("your-access-token")

	projects, err := client.GetProjects(context.Background())
	if err != nil {
		// handle error
		return
	}

	for _, p := range projects {
		fmt.Println(p.Name)
	}
}

func ExampleNewClient_withOptions() {
	client := ticktick.NewClient("your-access-token",
		ticktick.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}),
		// Use Dida365 API (TickTick's Chinese version).
		ticktick.WithBaseURL("https://api.dida365.com"),
	)

	projects, err := client.GetProjects(context.Background())
	if err != nil {
		// handle error
		return
	}

	for _, p := range projects {
		fmt.Println(p.Name)
	}
}

func ExampleClient_CreateTask() {
	client := ticktick.NewClient("your-access-token")

	task, err := client.CreateTask(context.Background(), &ticktick.CreateTaskRequest{
		Title:     "Buy groceries",
		ProjectID: "project-id",
		Priority:  ticktick.Int(ticktick.PriorityHigh),
	})
	if err != nil {
		// handle error
		return
	}

	fmt.Println(task.ID)
}

func ExampleClient_CreateTask_withDetails() {
	client := ticktick.NewClient("your-access-token")

	task, err := client.CreateTask(context.Background(), &ticktick.CreateTaskRequest{
		Title:     "Weekly report",
		ProjectID: "project-id",
		Content:   ticktick.String("Prepare and send the weekly status report"),
		Desc:      ticktick.String("Include metrics from dashboard"),
		IsAllDay:  ticktick.Bool(false),
		StartDate: ticktick.NewTime(time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)),
		DueDate:   ticktick.NewTime(time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)),
		TimeZone:  ticktick.String("America/New_York"),
		// Reminders use iCalendar TRIGGER format.
		// "TRIGGER:PT0S" — at the time of the event.
		// "TRIGGER:P0DT9H0M0S" — 9 hours before.
		// "TRIGGER:-PT30M" — 30 minutes before.
		Reminders: []string{"TRIGGER:PT0S", "TRIGGER:-PT30M"},
		// RepeatFlag uses iCalendar RRULE format.
		// "RRULE:FREQ=DAILY;INTERVAL=1" — every day.
		// "RRULE:FREQ=WEEKLY;INTERVAL=2" — every 2 weeks.
		// "RRULE:FREQ=MONTHLY;INTERVAL=1;BYDAY=1MO" — first Monday of each month.
		RepeatFlag: ticktick.String("RRULE:FREQ=WEEKLY;INTERVAL=1"),
		Priority:   ticktick.Int(ticktick.PriorityMedium),
		SortOrder:  ticktick.Int64(100),
		Items: []ticktick.CreateChecklistItemRequest{
			{Title: "Gather metrics"},
			{Title: "Write summary"},
			{Title: "Send to team"},
		},
	})
	if err != nil {
		// handle error
		return
	}

	fmt.Println(task.ID)
}

func ExampleClient_UpdateTask() {
	client := ticktick.NewClient("your-access-token")

	task, err := client.UpdateTask(context.Background(), "task-id", &ticktick.UpdateTaskRequest{
		ID:        "task-id",
		ProjectID: "project-id",
		Title:     ticktick.String("Updated title"),
		Priority:  ticktick.Int(ticktick.PriorityLow),
		DueDate:   ticktick.NewTime(time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		// handle error
		return
	}

	fmt.Println(task.Title)
}

func ExampleClient_CreateProject() {
	client := ticktick.NewClient("your-access-token")

	project, err := client.CreateProject(context.Background(), &ticktick.CreateProjectRequest{
		Name:     "Work Tasks",
		Color:    ticktick.String("#F18181"),
		ViewMode: ticktick.String("list"),
		Kind:     ticktick.String("TASK"),
	})
	if err != nil {
		// handle error
		return
	}

	fmt.Println(project.ID)
}

func ExampleClient_GetProjectData() {
	client := ticktick.NewClient("your-access-token")

	data, err := client.GetProjectData(context.Background(), "project-id")
	if err != nil {
		// handle error
		return
	}

	fmt.Printf("Project: %s\n", data.Project.Name)

	for _, task := range data.Tasks {
		fmt.Printf("  Task: %s (priority=%d)\n", task.Title, task.Priority)
	}

	for _, col := range data.Columns {
		fmt.Printf("  Column: %s\n", col.Name)
	}
}

func ExampleClient_GetTask_errorHandling() {
	client := ticktick.NewClient("your-access-token")

	task, err := client.GetTask(context.Background(), "project-id", "task-id")
	if err != nil {
		var apiErr *ticktick.Error
		if errors.As(err, &apiErr) {
			fmt.Printf("API error: HTTP %d: %s\n", apiErr.StatusCode, apiErr.Body)

			return
		}

		// Network or other error.
		fmt.Printf("unexpected error: %v\n", err)

		return
	}

	fmt.Println(task.Title)
}
