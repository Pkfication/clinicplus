package db

import (
	"log"
	"time"

	"clinicplus/internal/shared/config"
	"clinicplus/internal/shared/observability"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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

	// Enable GORM logging
	DB.LogMode(true)

	// Add callback to record metrics for database operations
	DB.Callback().Query().After("gorm:query").Register("observability:query", func(scope *gorm.Scope) {
		start := time.Now()
		duration := time.Since(start)
		observability.RecordDBQuery("query", scope.TableName(), duration, scope.DB().Error)
	})

	DB.Callback().Create().After("gorm:create").Register("observability:create", func(scope *gorm.Scope) {
		start := time.Now()
		duration := time.Since(start)
		observability.RecordDBQuery("create", scope.TableName(), duration, scope.DB().Error)
	})

	DB.Callback().Update().After("gorm:update").Register("observability:update", func(scope *gorm.Scope) {
		start := time.Now()
		duration := time.Since(start)
		observability.RecordDBQuery("update", scope.TableName(), duration, scope.DB().Error)
	})

	DB.Callback().Delete().After("gorm:delete").Register("observability:delete", func(scope *gorm.Scope) {
		start := time.Now()
		duration := time.Since(start)
		observability.RecordDBQuery("delete", scope.TableName(), duration, scope.DB().Error)
	})

	log.Println("Connected to the database with observability")
}
