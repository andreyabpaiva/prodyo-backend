package models

import (
	"time"

	"github.com/google/uuid"
)

// Indicator represents the productivity metrics for an iteration
// It contains calculated values for speed, rework, and instability
// The productivity levels are determined by project-level ranges
type Indicator struct {
	ID          uuid.UUID `json:"id"`
	IterationID uuid.UUID `json:"iteration_id"`

	// Calculated metric values (stored in DB for performance)
	SpeedValue       float64 `json:"speed_value"`       // tasks completed / time (in days)
	ReworkValue      float64 `json:"rework_value"`      // bugs / tasks
	InstabilityValue float64 `json:"instability_value"` // improvements / tasks

	// Calculated productivity levels based on project-level ranges
	// These are computed at runtime, not stored in DB
	SpeedLevel       ProductivityEnum `json:"speed_level,omitempty"`
	ReworkLevel      ProductivityEnum `json:"rework_level,omitempty"`
	InstabilityLevel ProductivityEnum `json:"instability_level,omitempty"`

	// Associated causes and actions for improvement
	Causes  []Cause  `json:"causes,omitempty"`
	Actions []Action `json:"actions,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CalculateProductivityLevels computes the productivity level for each indicator
// based on the provided project-level ranges
func (i *Indicator) CalculateProductivityLevels(ranges []IndicatorRange) {
	for _, r := range ranges {
		switch r.IndicatorType {
		case IndicatorSpeedPerIteration:
			i.SpeedLevel = r.Range.ClassifyValue(i.SpeedValue)
		case IndicatorReworkPerIteration:
			i.ReworkLevel = r.Range.ClassifyValue(i.ReworkValue)
		case IndicatorInstabilityIndex:
			i.InstabilityLevel = r.Range.ClassifyValue(i.InstabilityValue)
		}
	}
}

// GetMetricSummary returns a summary of all indicators with their values and levels
func (i *Indicator) GetMetricSummary() []IndicatorMetricValue {
	return []IndicatorMetricValue{
		{
			IndicatorType:     IndicatorSpeedPerIteration,
			Value:             i.SpeedValue,
			ProductivityLevel: i.SpeedLevel,
		},
		{
			IndicatorType:     IndicatorReworkPerIteration,
			Value:             i.ReworkValue,
			ProductivityLevel: i.ReworkLevel,
		},
		{
			IndicatorType:     IndicatorInstabilityIndex,
			Value:             i.InstabilityValue,
			ProductivityLevel: i.InstabilityLevel,
		},
	}
}
