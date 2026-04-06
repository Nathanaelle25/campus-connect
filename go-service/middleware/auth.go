package middleware

import (
	"net/http"
	"os"
)

func APIKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			expectedKey = "mysecretapikey"
		}

		if key != expectedKey {
			http.Error(w, `{"error": "Unauthorized / Invalid API Key"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
