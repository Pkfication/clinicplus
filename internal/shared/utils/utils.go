package utils

import (
	"database/sql"
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

type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time) // Marshal the valid time
	}
	return []byte("null"), nil // Marshal to null if invalid
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}
