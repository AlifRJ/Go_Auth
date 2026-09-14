package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/go-chi/jwtauth/v5"
)

// CheckTokenBlacklist memeriksa apakah JTI dari JWT Access Token ada di blacklist Redis
func CheckTokenBlacklist(authRepo model.AuthRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ambil claims dari JWT context
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: invalid context token",
				})
				return
			}

			// 2. Ekstrak JTI (JWT ID) dari claims
			jti, ok := claims["jti"].(string)
			if !ok || jti == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: missing token identifier (jti)",
				})
				return
			}

			// 3. Cek status blacklist di Redis melalui AuthRepository
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

			// 4. Lanjutkan ke handler berikutnya jika token aman
			next.ServeHTTP(w, r)
		})
	}
}

// Helper internal untuk mengembalikan response JSON
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}