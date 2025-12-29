package models

import (
	"time"

	"github.com/google/uuid"
)

type Action struct {
	ID          uuid.UUID `json:"id"`
	IndicatorRangeID uuid.UUID `json:"indicator_range_id"`
	Description string    `json:"description"`
	Cause       Cause     `json:"cause"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Assignee    User      `json:"assignee"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
