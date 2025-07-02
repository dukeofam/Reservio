package utils

import (
	"log"
	"net/http"
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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	_ = session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
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
	_, _ = w.Write([]byte(message))
}
