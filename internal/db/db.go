package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Init(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := migrate(db); err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	return db
}
