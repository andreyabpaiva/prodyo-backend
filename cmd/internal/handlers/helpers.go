package handlers

import (
	"prodyo-backend/cmd/internal/models"
	"strings"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	// Try RFC3339 format first
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// Try common formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, err
}

func normalizeStatus(statusStr string) models.StatusEnum {
	if statusStr == "" {
		return models.StatusNotStarted
	}

	lower := strings.ToLower(statusStr)
	lower = strings.ReplaceAll(lower, "_", "")
	lower = strings.ReplaceAll(lower, "-", "")

	switch lower {
	case "notstarted":
		return models.StatusNotStarted
	case "inprogress":
		return models.StatusInProgress
	case "completed":
		return models.StatusCompleted
	default:
		status := models.StatusEnum(statusStr)
		if status == models.StatusNotStarted || status == models.StatusInProgress || status == models.StatusCompleted {
			return status
		}
		return models.StatusNotStarted
	}
}
