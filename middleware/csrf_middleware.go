package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"reservio/config"
	"reservio/utils"
	"time"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, "session")
		token, _ := session.Values["csrf_token"].(string)
		expiry, _ := session.Values["csrf_token_expiry"].(int64)
		now := time.Now().Unix()

		// Log session state for debugging
		if os.Getenv("TEST_MODE") == "1" {
			log.Printf("[CSRF] Session values: %#v", session.Values)
			log.Printf("[CSRF] Current token: %s, expiry: %d, now: %d", token, expiry, now)
		}

		if token == "" || expiry == 0 || now > expiry {
			token = generateCSRFToken()
			session.Values["csrf_token"] = token
			session.Values["csrf_token_expiry"] = now + 7200 // 2 hours
			_ = session.Save(r, w)
			if os.Getenv("TEST_MODE") == "1" {
				log.Printf("[CSRF] Generated new token: %s", token)
			}
		}

		// Always set CSRF token in header in test mode
		if os.Getenv("TEST_MODE") == "1" {
			w.Header().Set("X-CSRF-Token", token)
			log.Printf("[CSRF] Set token in header: %s", token)
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			requestToken := r.Header.Get("X-CSRF-Token")
			if requestToken == "" {
				err := r.ParseForm()
				if err == nil {
					requestToken = r.FormValue("csrf_token")
				}
			}

			if os.Getenv("TEST_MODE") == "1" {
				log.Printf("[CSRF] Validating token - request: %s, session: %s", requestToken, token)
			}

			if requestToken != token {
				log.Printf("[CSRF] Invalid CSRF token: got=%s expected=%s", requestToken, token)
				utils.RespondWithValidationError(w, http.StatusForbidden, utils.NewValidationError("CSRF_INVALID", "Invalid CSRF token", nil))
				return
			}

			if os.Getenv("TEST_MODE") == "1" {
				log.Printf("[CSRF] Token validation successful")
			}
		}
		next.ServeHTTP(w, r)
	})
}

func RegenerateCSRFToken(w http.ResponseWriter, r *http.Request) error {
	session, _ := config.Store.Get(r, "session")
	token := generateCSRFToken()
	session.Values["csrf_token"] = token
	session.Values["csrf_token_expiry"] = time.Now().Unix() + 7200 // 2 hours
	if err := session.Save(r, w); err != nil {
		log.Printf("[CSRF] sess.Save error: %v", err)
		return err
	}
	// Always expose the fresh token so the frontend can store it
	w.Header().Set("X-CSRF-Token", token)
	return nil
}
