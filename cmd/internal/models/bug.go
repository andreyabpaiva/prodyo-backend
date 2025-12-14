package models

import (
	"time"

	"github.com/google/uuid"
)

type Bug struct {
	ID          uuid.UUID `json:"id"`
	TaskID      uuid.UUID `json:"task_id"`
	Assignee    User      `json:"assignee"`
	Number      int       `json:"number"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

