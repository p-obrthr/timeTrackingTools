package main

import (
	"database/sql"
	"fmt"
	"time"
)

type TimeLogDb struct {
	db *sql.DB
}

const file string = "timeTrackingTools.db"

const createDb string = `
	  CREATE TABLE IF NOT EXISTS timelogs (
		id INTEGER NOT NULL PRIMARY KEY,
		day INTEGER NOT NULL,
		month INTEGER NOT NULL,
		year INTEGER NOT NULL,
		week TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		kind INTEGER NOT NULL
	  );`

func InitDb() (*TimeLogDb, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createDb); err != nil {
		return nil, err
	}
	return &TimeLogDb{
		db: db,
	}, nil
}

func (db *TimeLogDb) Insert(timelog TimeLog) (int, error) {
	res, err := db.db.Exec("INSERT INTO timelogs (day, month, year, week, timestamp, kind) VALUES(?,?,?,?,?,?);", timelog.day, timelog.month, timelog.year, timelog.week, timelog.timestamp, timelog.kind)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *TimeLogDb) InsertDummyData() error {
	now := time.Now()
	year := now.Year()

	currentWeek := GetCalendarWeek()
	weeks := []int{currentWeek, currentWeek - 1, currentWeek - 2}

	for _, week := range weeks {
		for i := 0; i < week; i++ {
			timestamp := fmt.Sprintf("%02d:%02d", 0, 0)

			kind := 1

			timelog := TimeLog{
				day:       i + 1,
				month:     int(now.Month()),
				year:      year,
				week:      week,
				timestamp: timestamp,
				kind:      kind,
			}

			_, err := db.Insert(timelog)
			if err != nil {
				return fmt.Errorf("err: %v", err)
			}
		}
	}

	return nil
}

func (db *TimeLogDb) GetTimeLog(id int) (TimeLog, error) {
	row := db.db.QueryRow("SELECT * FROM timelogs WHERE id=?", id)

	timeLog := TimeLog{}
	if err := row.Scan(&timeLog.id, &timeLog.day, &timeLog.month, &timeLog.year, &timeLog.kind); err == sql.ErrNoRows {
		return TimeLog{}, err
	}
	return timeLog, nil
}

func (db *TimeLogDb) GetAllTimeLogs() ([]TimeLog, error) {
	rows, err := db.db.Query("SELECT * FROM timelogs;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanRows(rows)
}

func (db *TimeLogDb) GetTimeLogsByWeek(week int) ([]TimeLog, error) {
	rows, err := db.db.Query("SELECT * FROM timelogs WHERE week=?", week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanRows(rows)
}

func (db *TimeLogDb) scanRows(rows *sql.Rows) ([]TimeLog, error) {
	var timeLogs []TimeLog
	for rows.Next() {
		var timeLog TimeLog
		if err := rows.Scan(&timeLog.id, &timeLog.day, &timeLog.month, &timeLog.year, &timeLog.week, &timeLog.timestamp, &timeLog.kind); err != nil {
			return nil, err
		}
		timeLogs = append(timeLogs, timeLog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timeLogs, nil
}
