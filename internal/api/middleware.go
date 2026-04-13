package api

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func APIKeyAuth(apiKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		var key string
		if k := r.Header.Get("X-API-Key"); k != "" {
			key = k
		} else if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			key = strings.TrimPrefix(auth, "Bearer ")
		}

		if subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
