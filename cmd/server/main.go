package main

import (
	"clinicplus/internal/shared/config"
	"clinicplus/internal/shared/middleware"
	"clinicplus/pkg/server"
	"fmt"
	"log"

	"net/http"
)

func main() {
	// Get port from environment
	port := config.GetServerPort()

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)

	r := server.NewServer()
	handler := middleware.SetupCORS(r)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
