package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"reservio/config"
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
		if token == "" || expiry == 0 || now > expiry {
			token = generateCSRFToken()
			session.Values["csrf_token"] = token
			session.Values["csrf_token_expiry"] = now + 7200 // 2 hours
			session.Save(r, w)
		}

		if os.Getenv("TEST_MODE") == "1" && r.Method == http.MethodGet {
			w.Header().Set("X-CSRF-Token", token)
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			requestToken := r.Header.Get("X-CSRF-Token")
			if requestToken == "" {
				err := r.ParseForm()
				if err == nil {
					requestToken = r.FormValue("csrf_token")
				}
			}
			if requestToken != token {
				log.Printf("[CSRF] Invalid CSRF token: got=%s expected=%s", requestToken, token)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "Invalid CSRF token"}`))
				return
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
	return nil
}
