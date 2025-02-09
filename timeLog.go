package main

import (
	"fmt"
	"time"
)

type TimeLog struct {
	id        int
	day       int
	month     int
	year      int
	week      int
	timestamp string
	kind      int
}

func NewTimeLog(day int, month int, year int, kind int) *TimeLog {
	week := GetCalendarWeek()
	return &TimeLog{
		day:   day,
		month: month,
		year:  year,
		week:  week,
		kind:  kind,
	}
}

func NewTimeLogNow(kind int) *TimeLog {
	now := time.Now()
	week := GetCalendarWeek()
	timestamp := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	return &TimeLog{
		day:       now.Day(),
		month:     int(now.Month()),
		year:      now.Year(),
		week:      week,
		kind:      kind,
		timestamp: timestamp,
	}
}

func GetCalendarWeek() int {
	now := time.Now()
	_, week := now.ISOWeek()
	return week
}
