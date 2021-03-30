package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDatabase - creates a new gorm DB connection to our postgres database
func NewDatabase() (*gorm.DB, error) {
	fmt.Println("Starting new database connection")

	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")

	// Creates the connection string for postgres, and disables ssl for this demo code
	conStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbTable, dbPass)

	// opens a gorm DB
	db, err := gorm.Open("postgres", conStr)
	if err != nil {
		return db, err
	}

	// Make sure the database is online/pingable
	if err := db.DB().Ping(); err != nil {
		return db, err
	}

	// If no errors in startup or running confirmation, return the db and nil
	return db, nil
}
