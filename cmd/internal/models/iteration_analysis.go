package models

import "github.com/google/uuid"

type IterationAnalysisResponse struct {
	IterationID uuid.UUID                        `json:"iterationId"`
	Analysis    map[string]IndicatorAnalysisData `json:"analysis"`
}

type IndicatorAnalysisData struct {
	IndicatorType string         `json:"indicatorType"`
	XAxis         AxisDefinition `json:"xAxis"`
	YAxis         AxisDefinition `json:"yAxis"`
	Points        []DataPoint    `json:"points"`
}

type AxisDefinition struct {
	Type  string `json:"type,omitempty"`
	Label string `json:"label"`
}

type DataPoint struct {
	X      int              `json:"x"`
	Y      float64          `json:"y"`
	Status ProductivityEnum `json:"status"`
}
