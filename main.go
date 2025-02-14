package main

import (
	"clinicplus/db"
	"clinicplus/routes"
	"clinicplus/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func setupCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	database := db.InitDatabase()
	defer database.Close()

	// Get JWT secret from environment
	jwtSecret := utils.GetJWTSecret()

	// Create router
	router := mux.NewRouter()

	// Register routes and pass JWT secret
	routes.SetDB(database)
	routes.SetJWTKey(jwtSecret)
	routes.RegisterRoutes(router)

	// Log registered routes
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Registered Route: %v %v", methods, path)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking routes: %v", err)
	}

	// Get port from environment
	port := utils.GetServerPort()

	// Setup CORS middleware
	handler := setupCORS(router)

	// Start server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)

	// Start the server
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
