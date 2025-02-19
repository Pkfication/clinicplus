// /pkg/cron/cron.go
package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Function to clean up old records
func cleanupOldRecords() {
	// Your cleanup logic here
	log.Println("Cleaning up old records...")
}

// Function to initialize cron jobs
func StartCronJobs() {
	c := cron.New()

	// Schedule the cleanup job to run every day at midnight
	_, err := c.AddFunc("* 0 1 * *", cleanupOldRecords)
	if err != nil {
		log.Fatalf("Error scheduling cleanup job: %v", err)
	}

	// Start the cron scheduler
	c.Start()
}
