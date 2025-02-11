package routes

import (
    "net/http"
    "encoding/json"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    status := map[string]string{
        "status": "healthy",
        "version": "1.0.0",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

// Add to RegisterRoutes in routes/routes.go
func RegisterRoutes(r *mux.Router) {
    // ... existing routes
    r.HandleFunc("/health", HealthCheck).Methods("GET")
}