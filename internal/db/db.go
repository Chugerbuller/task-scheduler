package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title CHAR(32) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat CHAR(128) NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
func (s *Storage) Close() {
	s.db.Close()
}
