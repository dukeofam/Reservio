package controllers

import (
	"reservio/config"
	"reservio/models"
	"testing"
)

func TestAdminEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "adminuser@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Create slot using helper
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}
}
