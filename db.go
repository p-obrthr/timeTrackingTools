package main

import (
	"database/sql"
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

func (c *TimeLogDb) Insert(timelog TimeLog) (int, error) {
	res, err := c.db.Exec("INSERT INTO timelogs (day, month, year, entry_type) VALUES(?,?,?,?);", timelog.day, timelog.month, timelog.year, timelog.kind)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *TimeLogDb) GetTimeLog(id int) (TimeLog, error) {

	row := c.db.QueryRow("SELECT * FROM timelogs WHERE id=?", id)

	timeLog := TimeLog{}
	var err error
	if err = row.Scan(&timeLog.id, &timeLog.day, &timeLog.month, &timeLog.year, &timeLog.kind); err == sql.ErrNoRows {
		return TimeLog{}, err
	}
	return timeLog, err
}
