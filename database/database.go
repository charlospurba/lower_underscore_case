package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var DB *sql.DB

func ConnectDB() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=charlos dbname=lower_underscore_case sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Error connecting to database:", err)
	}

	// Tes koneksi database
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	fmt.Println("✅ Database Connected!")
}
