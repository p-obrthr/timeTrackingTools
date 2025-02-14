package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
	"timeTrackingTools/models"
)

type TimeLogDb struct {
	db *sql.DB
}

const dbDir string = "db"
const file string = "timeTrackingTools.db"

const createDb string = `
	CREATE TABLE IF NOT EXISTS timelogs (
		id INTEGER NOT NULL PRIMARY KEY,
		timestamp TEXT NOT NULL,
		week INTEGER NOT NULL,
		kind INTEGER NOT NULL
	);`

func InitDb() (*TimeLogDb, error) {
	dbPath := filepath.Join(dbDir, file)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createDb); err != nil {
		return nil, err
	}

	timeLogDb := &TimeLogDb{db: db}

	all, err := timeLogDb.GetAllTimeLogs()
	if err != nil {
		return nil, err
	}

	if len(all) == 0 {
		if err := timeLogDb.InsertDummyData(); err != nil {
			return nil, err
		}
	}

	return timeLogDb, nil
}

func (db *TimeLogDb) Insert(timelog models.TimeLog) (int, error) {

	formattedTimestamp := timelog.Timestamp.Format("2006-01-02 15:04:05")
	res, err := db.db.Exec("INSERT INTO timelogs (timestamp, week, kind) VALUES(?,?,?);", formattedTimestamp, timelog.Week, timelog.Kind)
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
	timelogs := []models.TimeLog{
		*models.NewTimeLog(time.Date(2025, 2, 10, 12, 0, 0, 0, time.Local), 0),
		*models.NewTimeLog(time.Date(2025, 2, 10, 14, 0, 0, 0, time.Local), 1),
		*models.NewTimeLog(time.Date(2025, 2, 9, 12, 0, 0, 0, time.Local), 0),
		*models.NewTimeLog(time.Date(2025, 2, 1, 12, 0, 0, 0, time.Local), 0),
	}

	for _, timelog := range timelogs {
		_, err := db.Insert(timelog)
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}

	}

	return nil
}

func (db *TimeLogDb) GetTimeLog(id int) (models.TimeLog, error) {
	row := db.db.QueryRow("SELECT * FROM timelogs WHERE id=?", id)

	timeLog := models.TimeLog{}
	if err := row.Scan(&timeLog.Id, &timeLog.Timestamp, &timeLog.Week, &timeLog.Kind); err == sql.ErrNoRows {
		return models.TimeLog{}, err
	}
	return timeLog, nil
}

func (db *TimeLogDb) GetAllTimeLogs() ([]models.TimeLog, error) {
	rows, err := db.db.Query("SELECT * FROM timelogs;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanRows(rows)
}

func (db *TimeLogDb) GetTimeLogsByWeek(week int) ([]models.TimeLog, error) {
	rows, err := db.db.Query("SELECT * FROM timelogs WHERE week=? ORDER BY timestamp ASC;", week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanRows(rows)
}

func (db *TimeLogDb) scanRows(rows *sql.Rows) ([]models.TimeLog, error) {
	var timeLogs []models.TimeLog
	for rows.Next() {
		var timeLog models.TimeLog
		var timestamp string

		if err := rows.Scan(&timeLog.Id, &timestamp, &timeLog.Week, &timeLog.Kind); err != nil {
			return nil, err
		}

		parsedTimestamp, _ := time.Parse("2006-01-02 15:04:05", timestamp)
		timeLog.Timestamp = parsedTimestamp.UTC()

		timeLogs = append(timeLogs, timeLog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return timeLogs, nil
}
