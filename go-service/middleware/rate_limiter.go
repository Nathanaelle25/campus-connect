package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu     sync.Mutex
	tokens map[string]int
	last   map[string]time.Time
}

var limiter = rateLimiter{
	tokens: make(map[string]int),
	last:   make(map[string]time.Time),
}

const maxTokens = 5
const refillRate = time.Second * 5

// RateLimitMiddleware applies a simple token bucket algorithm per IP strategy
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		if _, exists := limiter.last[ip]; !exists {
			limiter.tokens[ip] = maxTokens
			limiter.last[ip] = time.Now()
		} else {
			elapsed := time.Since(limiter.last[ip])
			if elapsed > refillRate {
				limiter.tokens[ip] = maxTokens
				limiter.last[ip] = time.Now()
			}
		}

		if limiter.tokens[ip] > 0 {
			limiter.tokens[ip]--
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, `{"error": "Rate limit exceeded"}`, http.StatusTooManyRequests)
		}
	})
}
