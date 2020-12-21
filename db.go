package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func getDB(filename string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("file:%s?", filename)
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	// Check table
	query := `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'reports';`
	var count int64
	err = db.Get(&count, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if count == 0 {
		query := `
			create table reports
			(
				id INTEGER primary key autoincrement,
				location text not null,
				icu_units INTEGER not null,
				inpatient_units INTEGER not null,
				created_at INTEGER not null
			);
		`
		_, err := db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("exec create table: %w", err)
		}
	}

	return db, nil
}

func saveReport(r *Report, db *sqlx.DB) error {
	query := "INSERT INTO reports (location, icu_units, inpatient_units, created_at) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, r.Location, r.ICUUnits, r.InpatientUnits, r.Time.UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}

func getLatestReport(db *sqlx.DB) (*Report, error) {
	var data struct {
		ID             int64  `db:"id"`
		Location       string `db:"location"`
		ICUUints       int64  `db:"icu_units"`
		InpatientUnits int64  `db:"inpatient_units"`
		T              int64  `db:"created_at"`
	}
	query := "SELECT * FROM reports ORDER BY id DESC LIMIT 1"
	err := db.Get(&data, query)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return &Report{
		Location:       data.Location,
		ICUUnits:       int(data.ICUUints),
		InpatientUnits: int(data.InpatientUnits),
		Time:           time.Unix(0, data.T),
	}, nil
}
