package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"reservio/config"
	"reservio/models"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

// Brute-force login attempt tracker (in-memory, can be replaced with Redis)
type LoginAttempt struct {
	Count      int
	LastFailed int64
}

var loginAttempts = struct {
	sync.Mutex
	m map[string]LoginAttempt
}{m: make(map[string]LoginAttempt)}

func IncrementLoginAttempt(email string) int {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	la := loginAttempts.m[email]
	la.Count++
	la.LastFailed = time.Now().Unix()
	loginAttempts.m[email] = la
	return la.Count
}

func ResetLoginAttempt(email string) {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	delete(loginAttempts.m, email)
}

func GetLoginAttempt(email string) LoginAttempt {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	return loginAttempts.m[email]
}

// Session helpers for gorilla/sessions
func SetSession(w http.ResponseWriter, r *http.Request, userID uint) {
	session, _ := config.Store.Get(r, "session")

	// Persist user_id and session_version
	session.Values["user_id"] = strconv.FormatUint(uint64(userID), 10)

	// Default session_version to 1; overwrite with DB value only if DB is available
	session.Values["session_version"] = "1"
	if config.DB != nil {
		var usr models.User
		if err := config.DB.Select("session_version").First(&usr, userID).Error; err == nil {
			session.Values["session_version"] = strconv.Itoa(usr.SessionVersion)
		}
	}

	// Ensure we have a CSRF token for this new/updated session
	token, _ := session.Values["csrf_token"].(string)
	if token == "" {
		// Generate a 32-byte random CSRF token (same entropy as middleware)
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		token = base64.StdEncoding.EncodeToString(b)
		session.Values["csrf_token"] = token
		session.Values["csrf_token_expiry"] = time.Now().Unix() + 7200 // 2 hours
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   os.Getenv("TEST_MODE") != "1",
		SameSite: http.SameSiteStrictMode,
	}
	if err := session.Save(r, w); err != nil {
		zap.L().Warn("SetSession save error", zap.Error(err))
	}

	// Expose CSRF token to the client so it can be stored
	w.Header().Set("X-CSRF-Token", token)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		zap.L().Warn("ClearSession save error", zap.Error(err))
	}
}

// Invalidate current session (global logout requires tracking all user sessions)
func InvalidateAllUserSessions(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, r)
}

// RespondWithError sends a JSON error response and logs the error
func RespondWithError(w http.ResponseWriter, code int, message string) {
	log.Printf("[ERROR] %s", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"error": message,
		"code":  ErrInvalidInput, // generic fallback; callers should prefer RespondWithValidationError for granular codes
	}

	jsonResp, _ := json.Marshal(response)
	_, _ = w.Write(jsonResp)
}

// RespondWithValidationError sends a JSON validation error response
func RespondWithValidationError(w http.ResponseWriter, code int, validationErr ValidationError) {
	log.Printf("[VALIDATION_ERROR] %s: %s", validationErr.Code, validationErr.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"error": validationErr.Message,
		"code":  validationErr.Code,
	}
	if validationErr.Details != nil {
		response["details"] = validationErr.Details
	}

	jsonResp, _ := json.Marshal(response)
	_, _ = w.Write(jsonResp)
}

// RespondWithSuccess sends a JSON success response
func RespondWithSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp, _ := json.Marshal(data)
	_, _ = w.Write(jsonResp)
}

// RespondWithPaginatedData sends a JSON response with paginated data
func RespondWithPaginatedData(w http.ResponseWriter, data interface{}, page, perPage, total int) {
	totalPages := (total + perPage - 1) / perPage // Ceiling division

	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	RespondWithSuccess(w, response)
}
