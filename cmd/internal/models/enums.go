package models

// ProductivityEnum represents the productivity level
type ProductivityEnum string

const (
	ProductivityOk       ProductivityEnum = "Ok"
	ProductivityAlert   ProductivityEnum = "Alert"
	ProductivityCritical ProductivityEnum = "Critical"
)

// MetricEnum represents the type of metric
type MetricEnum string

const (
	MetricWorkVelocity    MetricEnum = "WorkVelocity"
	MetricReworkIndex     MetricEnum = "ReworkIndex"
	MetricInstabilityIndex MetricEnum = "InstabilityIndex"
)

// StatusEnum represents the status of a task
type StatusEnum string

const (
	StatusNotStarted StatusEnum = "NotStarted"
	StatusInProgress StatusEnum = "InProgress"
	StatusCompleted  StatusEnum = "Completed"
)

