package models

import (
	"time"

	"github.com/google/uuid"
)

type Cause struct {
	ID                uuid.UUID        `json:"id"`
	IndicatorRangeID uuid.UUID        `json:"indicator_range_id"`
	Metric            MetricEnum       `json:"metric"`
	Description       string           `json:"description"`
	ProductivityLevel ProductivityEnum `json:"productivity_level"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

