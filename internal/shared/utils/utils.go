package utils

import (
	"encoding/json"
	"net/http"
)

// StandardResponse represents a consistent API response structure
type StandardResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
	Meta  interface{} `json:"meta"`
}

// SendJSONResponse sends a standardized JSON response
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}, err interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := StandardResponse{
		Data:  data,
		Error: err,
		Meta:  meta,
	}

	json.NewEncoder(w).Encode(response)
}
