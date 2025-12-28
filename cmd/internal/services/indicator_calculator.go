package services

import (
	"sort"
	"time"

	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
)

type IndicatorCalculator struct {
	tasks  []models.Task
	ranges map[models.IndicatorEnum]models.IndicatorRange
}

func NewIndicatorCalculator(tasks []models.Task, ranges []models.IndicatorRange) *IndicatorCalculator {
	rangeMap := make(map[models.IndicatorEnum]models.IndicatorRange)
	for _, r := range ranges {
		rangeMap[r.IndicatorType] = r
	}

	return &IndicatorCalculator{
		tasks:  tasks,
		ranges: rangeMap,
	}
}

func (ic *IndicatorCalculator) CalculateIterationAnalysis(iterationID uuid.UUID) models.IterationAnalysisResponse {
	completedTasks := ic.getCompletedTasksSorted()

	analysis := models.IterationAnalysisResponse{
		IterationID: iterationID,
		Analysis:    make(map[string]models.IndicatorAnalysisData),
	}

	analysis.Analysis["SpeedPerIteration"] = ic.calculateSpeedAnalysis(completedTasks)
	analysis.Analysis["ReworkPerIteration"] = ic.calculateReworkAnalysis(completedTasks)
	analysis.Analysis["InstabilityIndex"] = ic.calculateInstabilityAnalysis(completedTasks)

	return analysis
}

func (ic *IndicatorCalculator) getCompletedTasksSorted() []models.Task {
	var completed []models.Task
	for _, task := range ic.tasks {
		if task.Status == models.StatusCompleted {
			completed = append(completed, task)
		}
	}

	sort.Slice(completed, func(i, j int) bool {
		return completed[i].UpdatedAt.Before(completed[j].UpdatedAt)
	})

	return completed
}

func (ic *IndicatorCalculator) calculateSpeedAnalysis(completedTasks []models.Task) models.IndicatorAnalysisData {
	points := make([]models.DataPoint, 0, len(completedTasks))
	indicatorRange := ic.ranges[models.IndicatorSpeedPerIteration]

	for i, task := range completedTasks {
		var speed float64
		if task.Timer > 0 {
			speed = float64(task.Points) / float64(task.Timer)
		} else {
			speed = 0
		}

		status := ic.determineStatus(speed, indicatorRange)

		points = append(points, models.DataPoint{
			X:      i + 1,
			Y:      speed,
			Status: status,
		})
	}

	return models.IndicatorAnalysisData{
		IndicatorType: string(models.IndicatorSpeedPerIteration),
		XAxis: models.AxisDefinition{
			Type:  "TASK_SEQUENCE",
			Label: "Tasks concluídas",
		},
		YAxis: models.AxisDefinition{
			Label: "Pontos / Tempo (dias)",
		},
		Points: points,
	}
}

func (ic *IndicatorCalculator) calculateReworkAnalysis(completedTasks []models.Task) models.IndicatorAnalysisData {
	points := make([]models.DataPoint, 0, len(completedTasks))
	indicatorRange := ic.ranges[models.IndicatorReworkPerIteration]

	for i, task := range completedTasks {
		var rework float64
		for _, bug := range task.Bugs {
			rework += float64(bug.Points)
		}

		status := ic.determineStatus(rework, indicatorRange)

		points = append(points, models.DataPoint{
			X:      i + 1,
			Y:      rework,
			Status: status,
		})
	}

	return models.IndicatorAnalysisData{
		IndicatorType: string(models.IndicatorReworkPerIteration),
		XAxis: models.AxisDefinition{
			Type:  "TASK_SEQUENCE",
			Label: "Tasks concluídas",
		},
		YAxis: models.AxisDefinition{
			Label: "Bugs / Task",
		},
		Points: points,
	}
}

func (ic *IndicatorCalculator) calculateInstabilityAnalysis(completedTasks []models.Task) models.IndicatorAnalysisData {
	points := make([]models.DataPoint, 0, len(completedTasks))
	indicatorRange := ic.ranges[models.IndicatorInstabilityIndex]

	for i, task := range completedTasks {
		var instability float64
		for _, improvement := range task.Improvements {
			instability += float64(improvement.Points)
		}

		status := ic.determineStatus(instability, indicatorRange)

		points = append(points, models.DataPoint{
			X:      i + 1,
			Y:      instability,
			Status: status,
		})
	}

	return models.IndicatorAnalysisData{
		IndicatorType: string(models.IndicatorInstabilityIndex),
		XAxis: models.AxisDefinition{
			Type:  "TASK_SEQUENCE",
			Label: "Tasks concluídas",
		},
		YAxis: models.AxisDefinition{
			Label: "Melhorias / Task",
		},
		Points: points,
	}
}

func (ic *IndicatorCalculator) determineStatus(value float64, indicatorRange models.IndicatorRange) models.ProductivityEnum {
	r := indicatorRange.Range

	if value >= r.Critical.Min && value <= r.Critical.Max {
		return models.ProductivityCritical
	}

	if value >= r.Alert.Min && value <= r.Alert.Max {
		return models.ProductivityAlert
	}

	if value >= r.Ok.Min && value <= r.Ok.Max {
		return models.ProductivityOk
	}

	if value < r.Ok.Min {
		return models.ProductivityOk
	}
	return models.ProductivityCritical
}

func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
