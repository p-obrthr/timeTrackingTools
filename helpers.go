package main

import (
	"time"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetCalendarWeekNow() int {
	now := time.Now()
	_, week := now.ISOWeek()
	return week
}

func GetCalendarWeek(time time.Time) int {
	_, week := time.ISOWeek()
	return week
}
