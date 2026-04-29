package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB is the global database connection
var DB *sql.DB

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() {
	var err error
	dsn := "host=localhost port=5432 user=allen password=feedbackpass dbname=feedbackdb sslmode=disable"
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to the database!")
}

