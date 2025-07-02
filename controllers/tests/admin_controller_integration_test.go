package controllers

import (
	"testing"
)

func TestAdminEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	registerAndLogin(server, "adminuser@example.com", "testpassword123", csrfToken, cookie)

	// Create slot using helper
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}
}
