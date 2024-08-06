package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "data/database")

	if err != nil {
		log.Fatalf("Couldn't connect to the SQLite database: %s\n", err)
	}

	return db
}
