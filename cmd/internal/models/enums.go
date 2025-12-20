package models

// ProductivityEnum represents the productivity level
type ProductivityEnum string

const (
	ProductivityOk       ProductivityEnum = "Ok"
	ProductivityAlert    ProductivityEnum = "Alert"
	ProductivityCritical ProductivityEnum = "Critical"
)

// MetricEnum represents the type of metric (legacy, use IndicatorEnum for new code)
type MetricEnum string

const (
	MetricWorkVelocity     MetricEnum = "WorkVelocity"
	MetricReworkIndex      MetricEnum = "ReworkIndex"
	MetricInstabilityIndex MetricEnum = "InstabilityIndex"
)

// IndicatorEnum represents the type of productivity indicator
// These are calculated per iteration and classified based on project-level ranges
type IndicatorEnum string

const (
	// IndicatorSpeedPerIteration measures tasks completed per time unit (tasks/days)
	IndicatorSpeedPerIteration IndicatorEnum = "SpeedPerIteration"
	// IndicatorReworkPerIteration measures bug ratio (bugs/tasks)
	IndicatorReworkPerIteration IndicatorEnum = "ReworkPerIteration"
	// IndicatorInstabilityIndex measures improvement ratio (improvements/tasks)
	IndicatorInstabilityIndex IndicatorEnum = "InstabilityIndex"
)

// AllIndicatorTypes returns all available indicator types
func AllIndicatorTypes() []IndicatorEnum {
	return []IndicatorEnum{
		IndicatorSpeedPerIteration,
		IndicatorReworkPerIteration,
		IndicatorInstabilityIndex,
	}
}

// StatusEnum represents the status of a task
type StatusEnum string

const (
	StatusNotStarted StatusEnum = "NotStarted"
	StatusInProgress StatusEnum = "InProgress"
	StatusCompleted  StatusEnum = "Completed"
)
