package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS logins (id INTEGER PRIMARY KEY, username TEXT, tStamp TEXT, uuid TEXT, ipAddr TEXT, lat TEXT, lon TEXT, radius TEXT)")
	statement.Exec()

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}