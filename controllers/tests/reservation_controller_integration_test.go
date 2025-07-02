package controllers

import (
	"testing"
)

func TestReservationEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	registerAndLogin(server, "resparent+1@example.com", "testpassword123", csrfToken, cookie)

	// Create slot as admin
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}

	// Create child
	childID := createChild(server, csrfToken, cookie, "TestChild", "2018-01-01")
	if childID == 0 {
		t.Fatalf("expected child ID to be non-zero")
	}

	// Make reservation
	reservationID := createReservation(server, csrfToken, cookie, slotID, childID)
	if reservationID == 0 {
		t.Fatalf("expected reservation ID to be non-zero")
	}
}
