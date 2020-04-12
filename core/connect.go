package core

import (
	"database/sql"
	"log"

	// Load the PostgreSQL driver
	_ "github.com/lib/pq"
)

// DB holds the database connection to Postgres
var DB *sql.DB

// ConnectDB establishes database connection with postgres
func ConnectDB(connStr string) {
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	DB = db
}
