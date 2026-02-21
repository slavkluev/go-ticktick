package ticktick_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/slavkluev/go-ticktick"
)

func TestPriorityConstants(t *testing.T) {
	if ticktick.PriorityNone != 0 {
		t.Errorf("expected PriorityNone=0, got %d", ticktick.PriorityNone)
	}

	if ticktick.PriorityLow != 1 {
		t.Errorf("expected PriorityLow=1, got %d", ticktick.PriorityLow)
	}

	if ticktick.PriorityMedium != 3 {
		t.Errorf("expected PriorityMedium=3, got %d", ticktick.PriorityMedium)
	}

	if ticktick.PriorityHigh != 5 {
		t.Errorf("expected PriorityHigh=5, got %d", ticktick.PriorityHigh)
	}
}

func TestTaskStatusConstants(t *testing.T) {
	if ticktick.TaskStatusNormal != 0 {
		t.Errorf("expected TaskStatusNormal=0, got %d", ticktick.TaskStatusNormal)
	}

	if ticktick.TaskStatusCompleted != 2 {
		t.Errorf("expected TaskStatusCompleted=2, got %d", ticktick.TaskStatusCompleted)
	}
}

func TestChecklistStatusConstants(t *testing.T) {
	if ticktick.ChecklistStatusNormal != 0 {
		t.Errorf("expected ChecklistStatusNormal=0, got %d", ticktick.ChecklistStatusNormal)
	}

	if ticktick.ChecklistStatusCompleted != 1 {
		t.Errorf("expected ChecklistStatusCompleted=1, got %d", ticktick.ChecklistStatusCompleted)
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	tm := time.Date(2019, 11, 13, 3, 0, 0, 0, time.UTC)
	tt := ticktick.Time{Time: tm}

	data, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `"2019-11-13T03:00:00+0000"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestTimeMarshalJSONZero(t *testing.T) {
	var tt ticktick.Time

	data, err := json.Marshal(tt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != `""` {
		t.Errorf("expected empty string, got %s", string(data))
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	var tt ticktick.Time

	err := json.Unmarshal([]byte(`"2019-11-13T03:00:00+0000"`), &tt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tt.Year() != 2019 || tt.Month() != 11 || tt.Day() != 13 {
		t.Errorf("unexpected date: %v", tt.Time)
	}
}

func TestTimeUnmarshalJSONEmpty(t *testing.T) {
	var tt ticktick.Time

	err := json.Unmarshal([]byte(`""`), &tt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !tt.IsZero() {
		t.Errorf("expected zero time, got %v", tt.Time)
	}
}

func TestTimeUnmarshalJSONNull(t *testing.T) {
	var tt ticktick.Time

	err := json.Unmarshal([]byte(`null`), &tt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !tt.IsZero() {
		t.Errorf("expected zero time, got %v", tt.Time)
	}
}

func TestNewTime(t *testing.T) {
	tm := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pt := ticktick.NewTime(tm)

	if pt == nil {
		t.Fatal("expected non-nil pointer")
	}

	if !pt.Time.Equal(tm) {
		t.Errorf("expected %v, got %v", tm, pt.Time)
	}
}
