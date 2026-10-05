package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/go-chi/jwtauth/v5"
)

// CheckTokenBlacklist Checks if There are JTI and Access Token in Blacklist
func CheckTokenBlacklist(authRepo model.AuthRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Claims From Context
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: invalid context token",
				})
				return
			}

			// Extract JTI From Claims
			jti, ok := claims["jti"].(string)
			if !ok || jti == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: missing token identifier (jti)",
				})
				return
			}

			// Check Blacklist Status
			blacklisted, err := authRepo.IsTokenBlacklisted(r.Context(), jti)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "Internal server error during auth verification",
				})
				return
			}

			if blacklisted {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Token has been revoked. Please log in again.",
				})
				return
			}

			// Pass to Next Handler
			next.ServeHTTP(w, r)
		})
	}
}

// JSON Response Helper 
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}