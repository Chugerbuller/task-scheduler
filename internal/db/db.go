package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SchedularDb struct {
	db *sql.DB
}
var schedularDb SchedularDb
var Install bool = false

func Init(dbFile string) error {

	var schema string = `CREATE TABLE scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title CHAR(32) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat CHAR(128) NOT NULL
	)`

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if Install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	schedularDb = SchedularDb{
		db: db,
	}
	return nil
}
