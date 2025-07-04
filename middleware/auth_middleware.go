package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"reservio/config"
	"reservio/models"
	"reservio/utils"

	"go.uber.org/zap"
)

// ContextKey is a type for context keys used in request context
// to avoid collisions.
type ContextKey string

const UserIDKey ContextKey = "user_id"

// Protected middleware checks for a valid session and user_id
func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session")
		zap.L().Debug("Protected session", zap.Any("values", session.Values))
		idStr, ok := session.Values["user_id"].(string)
		if !ok || idStr == "" {
			zap.L().Debug("Protected no user_id", zap.Any("values", session.Values))
			if strings.HasPrefix(r.URL.Path, "/api/admin/") {
				utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "Resource not found", nil))
			} else {
				utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not logged in", nil))
			}
			return
		}
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			zap.L().Debug("Protected invalid user_id", zap.String("id", idStr))
			utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Invalid user ID", nil))
			return
		}
		id := uint(id64)
		// Compare session_version
		svStr, _ := session.Values["session_version"].(string)
		var usr models.User
		if err := config.DB.Select("session_version").First(&usr, id).Error; err == nil {
			if svStr == "" || svStr != strconv.Itoa(usr.SessionVersion) {
				zap.L().Debug("Session version mismatch", zap.String("cookie", svStr), zap.Int("db", usr.SessionVersion))
				utils.ClearSession(w, r)
				utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Session expired, please log in again", nil))
				return
			}
		}
		zap.L().Debug("Authenticated", zap.Uint("user_id", id))
		ctx := context.WithValue(r.Context(), UserIDKey, id)
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
			zap.L().Debug("AdminOnly user not found", zap.Uint("user_id", userID))
			utils.RespondWithValidationError(w, http.StatusForbidden, utils.NewValidationError(utils.ErrForbidden, "Forbidden", nil))
			return
		}
		if user.Role != "admin" {
			zap.L().Debug("AdminOnly forbidden role", zap.Uint("user_id", userID), zap.String("role", user.Role))
			utils.RespondWithValidationError(w, http.StatusForbidden, utils.NewValidationError(utils.ErrForbidden, "Forbidden", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}
