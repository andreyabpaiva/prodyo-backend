package models

import (
	"time"

	"github.com/google/uuid"
)

type Indicator struct {
	ID          uuid.UUID `json:"id"`
	IterationID uuid.UUID `json:"iteration_id"`
	Causes      []Cause   `json:"causes,omitempty"`
	Actions     []Action  `json:"actions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

