package routes

import (
	"clinicplus/internal/shared/utils"
	"net/http"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthData := map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	}

	utils.SendJSONResponse(w, http.StatusOK, healthData, nil, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"server":    "clinicplus-api",
	})
}
