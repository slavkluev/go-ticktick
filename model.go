package ticktick

// Task represents a TickTick task.
type Task struct {
	ID            string          `json:"id"`
	ProjectID     string          `json:"projectId"`
	Title         string          `json:"title"`
	IsAllDay      bool            `json:"isAllDay"`
	CompletedTime string          `json:"completedTime"`
	Content       string          `json:"content"`
	Desc          string          `json:"desc"`
	DueDate       string          `json:"dueDate"`
	Items         []ChecklistItem `json:"items"`
	Priority      int             `json:"priority"`
	Reminders     []string        `json:"reminders"`
	RepeatFlag    string          `json:"repeatFlag"`
	SortOrder     int64           `json:"sortOrder"`
	StartDate     string          `json:"startDate"`
	Status        int             `json:"status"`
	TimeZone      string          `json:"timeZone"`
	Kind          string          `json:"kind"`
}

// ChecklistItem represents a subtask within a task.
type ChecklistItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	CompletedTime string `json:"completedTime"`
	IsAllDay      bool   `json:"isAllDay"`
	SortOrder     int64  `json:"sortOrder"`
	StartDate     string `json:"startDate"`
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
	StartDate  *string                      `json:"startDate,omitempty"`
	DueDate    *string                      `json:"dueDate,omitempty"`
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
	StartDate  *string                      `json:"startDate,omitempty"`
	DueDate    *string                      `json:"dueDate,omitempty"`
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
	StartDate     *string `json:"startDate,omitempty"`
	IsAllDay      *bool   `json:"isAllDay,omitempty"`
	SortOrder     *int64  `json:"sortOrder,omitempty"`
	TimeZone      *string `json:"timeZone,omitempty"`
	Status        *int    `json:"status,omitempty"`
	CompletedTime *string `json:"completedTime,omitempty"`
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

// String returns a pointer to the given string value.
func String(v string) *string { return &v }

// Int returns a pointer to the given int value.
func Int(v int) *int { return &v }

// Int64 returns a pointer to the given int64 value.
func Int64(v int64) *int64 { return &v }

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }
