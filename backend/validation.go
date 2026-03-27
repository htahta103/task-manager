package main

import (
	"strings"
	"time"
)

func normalizeTitle(s string) string {
	return strings.TrimSpace(s)
}

func validateTitle(title string) bool {
	title = normalizeTitle(title)
	return title != "" && len(title) <= 255
}

func isValidStatus(s TaskStatus) bool {
	switch s {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

func isValidPriority(p TaskPriority) bool {
	switch p {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh:
		return true
	default:
		return false
	}
}

func parseISODatePtr(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	if strings.TrimSpace(*s) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &utc, nil
}

func formatISODatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format("2006-01-02")
	return &s
}
