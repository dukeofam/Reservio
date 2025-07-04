package middleware

import (
	"net/http"
	"os"
	"strings"
)

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- Security headers -------------------------------------------------
		// Only enable HSTS when request is HTTPS (r.TLS != nil) OR explicitly enabled via env.
		if r.TLS != nil || os.Getenv("ENABLE_HSTS") == "1" {
			// 2-year max-age recommended by OWASP
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content-Security-Policy can be overridden via env (single string)
		csp := os.Getenv("CSP")
		if csp == "" {
			csp = "default-src 'self'; frame-ancestors 'none'; base-uri 'self'"
		}
		w.Header().Set("Content-Security-Policy", csp)

		// Disable features we don't use
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		// ---------------------------------------------------------------------

		// CORS headers
		origin := r.Header.Get("Origin")
		if origin != "" {
			allowed := false

			// Test mode – always allow localhost dev server origins
			if os.Getenv("TEST_MODE") == "1" && (origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000") {
				allowed = true
			}

			// Environment variable override (comma-separated list)
			if !allowed {
				if ao := os.Getenv("ALLOWED_ORIGINS"); ao != "" {
					origins := strings.Split(ao, ",")
					for _, o := range origins {
						if strings.TrimSpace(o) == origin {
							allowed = true
							break
						}
					}
				}
			}

			// Default production fallback (if ENV var not set)
			if !allowed && os.Getenv("ENVIRONMENT") == "production" && origin == "https://yourdomain.com" {
				allowed = true
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			}
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
