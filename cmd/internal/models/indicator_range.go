package models

import (
	"time"

	"github.com/google/uuid"
)

// IndicatorRange represents the productivity ranges for a specific indicator type within a project
// Each project has separate ranges for SpeedPerIteration, ReworkPerIteration, and InstabilityIndex
// These ranges are set at project creation and apply to all iterations
type IndicatorRange struct {
	ID            uuid.UUID         `json:"id"`
	ProjectID     uuid.UUID         `json:"project_id"`
	IndicatorType IndicatorEnum     `json:"indicator_type"`
	Range         ProductivityRange `json:"range"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// IndicatorMetricValue represents the calculated value for a specific indicator
// along with its classification based on the configured project-level ranges
type IndicatorMetricValue struct {
	IndicatorType     IndicatorEnum    `json:"indicator_type"`
	Value             float64          `json:"value"`
	ProductivityLevel ProductivityEnum `json:"productivity_level"`
}

// CalculateProductivityLevel determines the productivity level for this indicator
// based on the provided range configuration
func (imv *IndicatorMetricValue) CalculateProductivityLevel(ranges ProductivityRange) ProductivityEnum {
	return ranges.ClassifyValue(imv.Value)
}
