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
	registerAndLogin(server, "adminuser@example.com", "testpassword123", csrfToken, cookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", "adminuser@example.com").Update("role", "admin")

	// Create slot using helper
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}
}
