package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter stores rate limiters per IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// GetLimiter returns the rate limiter for the provided IP address
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// Cleanup removes old limiters (call periodically)
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// In a production app, you'd implement cleanup logic here
	// For now, we'll keep all limiters in memory
}

// Global rate limiter instance
var globalLimiter = NewRateLimiter(5, 20) // 5 requests per second, burst of 20

// RateLimitMiddleware applies rate limiting per IP address
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting in test mode or CI
		if os.Getenv("TEST_MODE") == "1" || os.Getenv("CI") == "true" {
			next.ServeHTTP(w, r)
			return
		}

		// Get client IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		// Get rate limiter for this IP
		limiter := globalLimiter.GetLimiter(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "Rate limit exceeded. Please try again later.",
			}); err != nil {
				zap.L().Error("encode error", zap.Error(err))
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
