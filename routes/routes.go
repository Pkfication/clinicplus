package routes

import (
    "github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
    // Existing route registrations
    r.HandleFunc("/login", Login).Methods("POST")
    r.HandleFunc("/logout", Logout).Methods("POST")
    // ... other existing routes

    // Add health check route
    r.HandleFunc("/health", HealthCheck).Methods("GET")
}