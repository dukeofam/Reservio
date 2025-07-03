package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"reservio/config"

	"github.com/gorilla/sessions"
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
	session.Values["user_id"] = strconv.Itoa(int(userID))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   os.Getenv("TEST_MODE") != "1",
		SameSite: http.SameSiteStrictMode,
	}
	if err := session.Save(r, w); err != nil {
		log.Printf("[SetSession] session.Save error: %v", err)
	}
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		log.Printf("[ClearSession] session.Save error: %v", err)
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
	jsonResp := []byte(`{"error": "` + message + `"}`)
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
