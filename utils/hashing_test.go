package utils

import (
	"net/http/httptest"
	"os"
	"reservio/config"
	"testing"
)

func TestMain(m *testing.M) {
	config.InitSessionStore()
	os.Exit(m.Run())
}

func TestSetSessionAndClearSession(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	SetSession(w, r, 42)

	req := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}
	w2 := httptest.NewRecorder()
	session, _ := config.Store.Get(req, "session")
	userID, _ := session.Values["user_id"].(string)
	if userID != "42" {
		t.Errorf("expected user_id 42, got %v", userID)
	}

	ClearSession(w2, req)
	req2 := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range w2.Result().Cookies() {
		req2.AddCookie(cookie)
	}
	session2, _ := config.Store.Get(req2, "session")
	if session2.Values["user_id"] != nil {
		t.Errorf("expected user_id to be cleared, got %v", session2.Values["user_id"])
	}
}
