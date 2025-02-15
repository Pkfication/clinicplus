package db

import (
	"log"

	"clinicplus/internal/shared/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDatabase() {
	// Load environment variables
	config.LoadEnv()

	// Get database connection string from environment
	dbURL := config.GetDatabaseURL()

	// Open database connection
	var err error
	DB, err = gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to the database")
}
