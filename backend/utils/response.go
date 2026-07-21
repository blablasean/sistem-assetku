package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendSuccess sends a success response
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, statusCode int, message string, err string) {
	response := Response{
		Status:  false,
		Message: message,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
