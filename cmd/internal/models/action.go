package models

import (
	"time"

	"github.com/google/uuid"
)

type Action struct {
	ID          uuid.UUID `json:"id"`
	IndicatorID uuid.UUID `json:"indicator_id"`
	Description string    `json:"description"`
	Cause       Cause     `json:"cause"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

