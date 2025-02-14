package db

import (
    "log"
    
   "clinicplus/utils"
   "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

func InitDatabase() *gorm.DB {
    // Load environment variables
    utils.LoadEnv()

    // Get database connection string from environment
    dbURL := utils.GetDatabaseURL()

    // Open database connection
    db, err := gorm.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    return db
}