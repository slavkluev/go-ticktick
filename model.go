package ticktick

import (
	"fmt"
	"time"
)

// TimeLayout is the time format used by the TickTick API.
const TimeLayout = "2006-01-02T15:04:05-0700"

// Task priority levels.
const (
	PriorityNone   = 0
	PriorityLow    = 1
	PriorityMedium = 3
	PriorityHigh   = 5
)

// Task completion status values.
const (
	TaskStatusNormal    = 0
	TaskStatusCompleted = 2
)

// ChecklistItem completion status values.
const (
	ChecklistStatusNormal    = 0
	ChecklistStatusCompleted = 1
)

// Task represents a TickTick task.
type Task struct {
	ID            string          `json:"id"`
	ProjectID     string          `json:"projectId"`
	Title         string          `json:"title"`
	IsAllDay      bool            `json:"isAllDay"`
	CompletedTime Time            `json:"completedTime"`
	Content       string          `json:"content"`
	Desc          string          `json:"desc"`
	DueDate       Time            `json:"dueDate"`
	Items         []ChecklistItem `json:"items"`
	Priority      int             `json:"priority"`
	Reminders     []string        `json:"reminders"`
	RepeatFlag    string          `json:"repeatFlag"`
	SortOrder     int64           `json:"sortOrder"`
	StartDate     Time            `json:"startDate"`
	Status        int             `json:"status"`
	TimeZone      string          `json:"timeZone"`
	Kind          string          `json:"kind"`
}

// ChecklistItem represents a subtask within a task.
type ChecklistItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	CompletedTime Time   `json:"completedTime"`
	IsAllDay      bool   `json:"isAllDay"`
	SortOrder     int64  `json:"sortOrder"`
	StartDate     Time   `json:"startDate"`
	TimeZone      string `json:"timeZone"`
}

// Project represents a TickTick project.
type Project struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	SortOrder  int64  `json:"sortOrder"`
	Closed     bool   `json:"closed"`
	GroupID    string `json:"groupId"`
	ViewMode   string `json:"viewMode"`
	Permission string `json:"permission"`
	Kind       string `json:"kind"`
}

// Column represents a kanban column within a project.
type Column struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	SortOrder int64  `json:"sortOrder"`
}

// ProjectData holds a project along with its tasks and columns.
type ProjectData struct {
	Project Project  `json:"project"`
	Tasks   []Task   `json:"tasks"`
	Columns []Column `json:"columns"`
}

// CreateTaskRequest contains the fields for creating a new task.
type CreateTaskRequest struct {
	Title      string                       `json:"title"`
	ProjectID  string                       `json:"projectId"`
	Content    *string                      `json:"content,omitempty"`
	Desc       *string                      `json:"desc,omitempty"`
	IsAllDay   *bool                        `json:"isAllDay,omitempty"`
	StartDate  *Time                        `json:"startDate,omitempty"`
	DueDate    *Time                        `json:"dueDate,omitempty"`
	TimeZone   *string                      `json:"timeZone,omitempty"`
	Reminders  []string                     `json:"reminders,omitempty"`
	RepeatFlag *string                      `json:"repeatFlag,omitempty"`
	Priority   *int                         `json:"priority,omitempty"`
	SortOrder  *int64                       `json:"sortOrder,omitempty"`
	Items      []CreateChecklistItemRequest `json:"items,omitempty"`
}

// UpdateTaskRequest contains the fields for updating an existing task.
type UpdateTaskRequest struct {
	ID         string                       `json:"id"`
	ProjectID  string                       `json:"projectId"`
	Title      *string                      `json:"title,omitempty"`
	Content    *string                      `json:"content,omitempty"`
	Desc       *string                      `json:"desc,omitempty"`
	IsAllDay   *bool                        `json:"isAllDay,omitempty"`
	StartDate  *Time                        `json:"startDate,omitempty"`
	DueDate    *Time                        `json:"dueDate,omitempty"`
	TimeZone   *string                      `json:"timeZone,omitempty"`
	Reminders  []string                     `json:"reminders,omitempty"`
	RepeatFlag *string                      `json:"repeatFlag,omitempty"`
	Priority   *int                         `json:"priority,omitempty"`
	SortOrder  *int64                       `json:"sortOrder,omitempty"`
	Items      []CreateChecklistItemRequest `json:"items,omitempty"`
}

// CreateChecklistItemRequest contains the fields for a subtask in a create or update request.
type CreateChecklistItemRequest struct {
	Title         string  `json:"title"`
	StartDate     *Time   `json:"startDate,omitempty"`
	IsAllDay      *bool   `json:"isAllDay,omitempty"`
	SortOrder     *int64  `json:"sortOrder,omitempty"`
	TimeZone      *string `json:"timeZone,omitempty"`
	Status        *int    `json:"status,omitempty"`
	CompletedTime *Time   `json:"completedTime,omitempty"`
}

// CreateProjectRequest contains the fields for creating a new project.
type CreateProjectRequest struct {
	Name      string  `json:"name"`
	Color     *string `json:"color,omitempty"`
	SortOrder *int64  `json:"sortOrder,omitempty"`
	ViewMode  *string `json:"viewMode,omitempty"`
	Kind      *string `json:"kind,omitempty"`
}

// UpdateProjectRequest contains the fields for updating an existing project.
type UpdateProjectRequest struct {
	Name      *string `json:"name,omitempty"`
	Color     *string `json:"color,omitempty"`
	SortOrder *int64  `json:"sortOrder,omitempty"`
	ViewMode  *string `json:"viewMode,omitempty"`
	Kind      *string `json:"kind,omitempty"`
}

// Time wraps [time.Time] with custom JSON marshaling for the TickTick API date format.
type Time struct {
	time.Time
}

// NewTime creates a pointer to a Time value. Useful for optional date fields in request types.
func NewTime(t time.Time) *Time {
	return &Time{Time: t}
}

// MarshalJSON serializes a Time to JSON using the TickTick date format.
// A zero Time marshals to an empty string.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}

	return []byte(`"` + t.Format(TimeLayout) + `"`), nil
}

// UnmarshalJSON deserializes a Time from JSON. Empty strings and null values
// are unmarshaled as zero time.
func (t *Time) UnmarshalJSON(data []byte) error {
	s := string(data)

	if s == "null" {
		return nil
	}

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "" {
		t.Time = time.Time{}

		return nil
	}

	parsed, err := time.Parse(TimeLayout, s)
	if err != nil {
		return fmt.Errorf("ticktick: cannot parse time %q: %w", s, err)
	}

	t.Time = parsed

	return nil
}

// String returns a pointer to the given string value.
func String(v string) *string { return &v }

// Int returns a pointer to the given int value.
func Int(v int) *int { return &v }

// Int64 returns a pointer to the given int64 value.
func Int64(v int64) *int64 { return &v }

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }
