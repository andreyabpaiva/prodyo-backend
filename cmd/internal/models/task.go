package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID           uuid.UUID  `json:"id"`
	IterationID  uuid.UUID  `json:"iteration_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Assignee     User       `json:"assignee"`
	Status       StatusEnum `json:"status"`
	Timer        int64      `json:"timer"`
	Points       int        `json:"points"`
	ExpectedTime float64    `json:"expected_time"` 
	Tasks        []Task     `json:"tasks,omitempty"` // Sub-tasks
	Improvements []Improv   `json:"improvements,omitempty"`
	Bugs         []Bug      `json:"bugs,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
