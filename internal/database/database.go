package database

import (
	"fmt"
	"os"
	"time"

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

	// Sometimes docker takes a while to get the db up and running - the dependency in docker only waits for the
	// service to be running, not ready. So, allow some retries.
	retries := 3

	for err != nil {
		fmt.Printf("Failed to connect to database - retries remaining: %d\n", retries)
		if retries > 0 {
			retries--
			time.Sleep(5 * time.Second)
			fmt.Println(conStr)
			db, err = gorm.Open("postgres", conStr)
			continue
		} else {
			break
		}
	}

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
