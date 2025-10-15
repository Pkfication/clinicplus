package main

import (
	"clinicplus/internal/shared/config"
	"clinicplus/internal/shared/middleware"
	"clinicplus/internal/shared/observability"
	"clinicplus/pkg/cron"
	"clinicplus/pkg/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"
)

func main() {
	// Initialize tracing
	cleanup := observability.InitTracing("clinicplus-api", "1.0.0")
	defer cleanup()

	// Start cron jobs
	cron.StartCronJobs()
	
	// Get port from environment
	port := config.GetServerPort()
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)

	// Create server with observability
	r := server.NewServer()
	handler := middleware.SetupCORS(r)

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", serverAddr)
		if err := http.ListenAndServe(serverAddr, handler); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
