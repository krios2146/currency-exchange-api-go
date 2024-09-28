package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "data/database")

	if err != nil {
		slog.Error("Couldn't connect to the SQLite database", "error", err)
		os.Exit(1)
	}

	return db
}
