package middleware

import (
	"encoding/json"
	"net/http"
)

// 1. Definisikan struktur JSON yang seragam dengan error API Anda yang lain
type APIErrorResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// 2. Buat fungsi Custom Rate Limit Error
func WriteRateLimitError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Retry-After", "60") 
	
	w.WriteHeader(http.StatusTooManyRequests)

	response := APIErrorResponse{
		Status:  http.StatusTooManyRequests,
		Error:   "Too Many Requests",
		Message: "You are sending request too fast. Please wait a few minutes.",
		Code:    "ERR_RATE_LIMIT_EXCEEDED",
	}

	_ = json.NewEncoder(w).Encode(response)
}
