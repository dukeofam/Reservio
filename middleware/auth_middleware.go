package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"reservio/config"
	"reservio/models"
)

// ContextKey is a type for context keys used in request context
// to avoid collisions.
type ContextKey string

const UserIDKey ContextKey = "user_id"

// Protected middleware checks for a valid session and user_id
func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session")
		idStr, ok := session.Values["user_id"].(string)
		if !ok || idStr == "" {
			http.Error(w, `{"error": "Not logged in"}`, http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, uint(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly middleware checks if the user is an admin
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(UserIDKey).(uint)
		if !ok {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			log.Printf("[AdminOnly] Forbidden: user not found (user_id=%d)", userID)
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		if user.Role != "admin" {
			log.Printf("[AdminOnly] Forbidden: user_id=%d, role=%s", userID, user.Role)
			http.Error(w, `{"error": "Forbidden"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
