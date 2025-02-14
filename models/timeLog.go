package models

import (
	"time"
	"timeTrackingTools/helpers"
)

type TimeLog struct {
	Id        int
	Timestamp time.Time
	Week      int
	Kind      int
}

func NewTimeLog(timestamp time.Time, kind int) *TimeLog {
	timestamp = time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), 0, 0, time.UTC)
	week := helpers.GetCalendarWeek(timestamp)

	return &TimeLog{
		Timestamp: timestamp,
		Week:      week,
		Kind:      kind,
	}
}

func NewTimeLogNow(kind int) *TimeLog {
	now := time.Now()
	week := helpers.GetCalendarWeek(now)
	timestamp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	return &TimeLog{
		Week:      week,
		Kind:      kind,
		Timestamp: timestamp,
	}
}
