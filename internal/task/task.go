package task

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Status int

const (
	StatusTodo       Status = 0 
	StatusInProgress Status = 1 
	StatusDone       Status = 2 
)

//status to string
func (s Status) String() string {
	switch s {
	case StatusInProgress:
		return "in-progress"
	case StatusDone:
		return "done"
	default:
		return "todo"
	}
}


//how the task will be represented on disk in json (tasks.json)
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Priority  int       `json:"priority"`  
	Status    Status    `json:"status"`
	Due       string    `json:"due"`       
	Tags      []string  `json:"tags"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}


//priority
func (t Task) PrioritySymbol() string {
	switch t.Priority {
	case 2:
		return "!!"
	case 3:
		return "!!!"
	default:
		return "!"
	}
}


//priority to int
func ParsePriority(s string) (int, error) {
	switch strings.TrimSpace(s) {
	case "!":
		return 1, nil
	case "!!":
		return 2, nil
	case "!!!":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid priority %q — use !, !!, or !!!", s)
	}
}



func (t Task) IsOverdue() bool {
	if t.Due == "" || t.Status == StatusDone {
		return false
	}
	due, err := time.Parse("2006-01-02", t.Due)
	if err != nil {
		return false
	}

	return time.Now().After(due.Add(24 * time.Hour))
}


func (t Task) IsDueToday() bool {
	if t.Due == "" || t.Status == StatusDone {
		return false
	}
	due, err := time.Parse("2006-01-02", t.Due)
	if err != nil {
		return false
	}
	now := time.Now()

	return due.Year() == now.Year() && due.YearDay() == now.YearDay()
}


func Filter(tasks []Task, tag, statusFilter string) []Task {

	if tag == "" && statusFilter == "" {
		return tasks
	}

	out := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if tag != "" && !hasTag(t, tag) {
			continue
		}

		if statusFilter != "" && t.Status.String() != statusFilter {
			continue
		}
		out = append(out, t)
	}
	return out
}


func hasTag(t Task, tag string) bool {
	for _, tg := range t.Tags {
		if strings.EqualFold(tg, tag) {
			return true
		}
	}
	return false
}


func Sort(tasks []Task, by string) []Task {

	sorted := make([]Task, len(tasks))
	copy(sorted, tasks)

	switch by {
	case "priority":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Priority > sorted[j].Priority
		})
	case "due":
		sort.Slice(sorted, func(i, j int) bool {
			a, b := sorted[i].Due, sorted[j].Due
			if a == "" {
				return false 
			}
			if b == "" {
				return true
			}
			return a < b 
		})
	case "created":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		})
	default: 
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ID < sorted[j].ID
		})
	}
	return sorted
}
