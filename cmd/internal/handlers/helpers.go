package handlers

import (
	"fmt"
	"prodyo-backend/cmd/internal/models"
	"strings"
	"time"
)

func parseDuration(durationStr string) (int64, error) {
	d, err := time.ParseDuration(durationStr)
	if err == nil {
		return int64(d.Seconds()), nil
	}

	var seconds int64
	_, err = fmt.Sscanf(durationStr, "%d", &seconds)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s (use format like '2h', '90m', or '7200' for seconds)", durationStr)
	}

	return seconds, nil
}

func parseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

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
