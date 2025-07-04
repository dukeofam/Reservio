package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"reservio/config"
	"reservio/models"
	"reservio/utils"
)

// ContextKey is a type for context keys used in request context
// to avoid collisions.
type ContextKey string

const UserIDKey ContextKey = "user_id"

// Protected middleware checks for a valid session and user_id
func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session")
		log.Printf("[Protected] session.Values: %#v", session.Values)
		idStr, ok := session.Values["user_id"].(string)
		if !ok || idStr == "" {
			log.Printf("[Protected] No user_id in session: %#v", session.Values)
			if strings.HasPrefix(r.URL.Path, "/api/admin/") {
				utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "Resource not found", nil))
			} else {
				utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not logged in", nil))
			}
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("[Protected] Invalid user_id in session: %v", idStr)
			utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Invalid user ID", nil))
			return
		}
		// Compare session_version
		svStr, _ := session.Values["session_version"].(string)
		var usr models.User
		if err := config.DB.Select("session_version").First(&usr, id).Error; err == nil {
			if svStr == "" || svStr != strconv.Itoa(usr.SessionVersion) {
				log.Printf("[Protected] Session version mismatch: cookie=%s db=%d", svStr, usr.SessionVersion)
				utils.ClearSession(w, r)
				utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Session expired, please log in again", nil))
				return
			}
		}
		log.Printf("[Protected] Authenticated user_id: %d", id)
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
			utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Unauthorized", nil))
			return
		}
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			log.Printf("[AdminOnly] Forbidden: user not found (user_id=%d)", userID)
			utils.RespondWithValidationError(w, http.StatusForbidden, utils.NewValidationError(utils.ErrForbidden, "Forbidden", nil))
			return
		}
		if user.Role != "admin" {
			log.Printf("[AdminOnly] Forbidden: user_id=%d, role=%s", userID, user.Role)
			utils.RespondWithValidationError(w, http.StatusForbidden, utils.NewValidationError(utils.ErrForbidden, "Forbidden", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}
