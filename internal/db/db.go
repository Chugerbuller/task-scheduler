package db

import (
	"TaskScheduler/internal/logger"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db   *sql.DB
	logs *logger.Logger
}

var storage *Storage

// The Init function initializes a SQLite database connection and creates a table for scheduling tasks
// if it does not already exist.
func Init(dbPath string, logger *logger.Logger) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title CHAR(32) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat CHAR(128) NOT NULL
	);
	CREATE INDEX  IF NOT EXISTS idx_date ON scheduler(date)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	storage = &Storage{
		db:   db,
		logs: logger,
	}
	storage.logs.Print("Connected.")
	return nil
}

// The Close function closes the database connection in Go.
func Close() {
	storage.db.Close()
}
