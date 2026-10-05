package middleware

import (
	"net"
	"net/http"
	"strings"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientIP string

		// Prioritize True-Client-IP
		if trueIP := r.Header.Get("True-Client-IP"); trueIP != "" {
			clientIP = strings.TrimSpace(trueIP)
		} 
		// Fallback to Parsing X-Forwarded-For Securely from the Right Side (1 hop)
		if clientIP == "" {
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				parts := strings.Split(xff, ",")
				if len(parts) > 0 {
					clientIP = strings.TrimSpace(parts[0]) 
				}
			}
		}

		// IP Clean Up
		if clientIP != "" {
			if host, _, err := net.SplitHostPort(clientIP); err == nil {
				r.RemoteAddr = host
			} else {
				r.RemoteAddr = clientIP
			}
		}

		next.ServeHTTP(w, r)
	})
}