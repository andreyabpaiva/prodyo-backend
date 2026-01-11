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
	Summary       *SpeedSummary  `json:"summary,omitempty"`
	Values        *SpeedValues   `json:"values,omitempty"`
}

type AxisDefinition struct {
	Type  string `json:"type,omitempty"`
	Label string `json:"label"`
}

type DataPoint struct {
	X      interface{}      `json:"x"`
	Y      float64          `json:"y"`
	Status ProductivityEnum `json:"status,omitempty"`
}

type SpeedSummary struct {
	TotalPoints        int     `json:"totalPoints"`
	TotalEstimatedTime float64 `json:"totalEstimatedTime"`
	TotalActualTime    float64 `json:"totalActualTime"`
}

type SpeedValues struct {
	ExpectedSpeed     float64 `json:"expectedSpeed"`
	ActualSpeed       float64 `json:"actualSpeed"`
	EfficiencyPercent float64 `json:"efficiencyPercent"`
}
