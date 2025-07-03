package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reservio/config"
	"reservio/models"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "adminuser@example.com"
	password := "testpassword123"
	csrfToken, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Create slot using helper
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}
}

func TestAdminSlotManagement(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "adminuser2@example.com"
	password := "testpassword123"
	csrfToken, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Test valid slot creation
	slotPayload := map[string]interface{}{"date": "2025-12-15", "capacity": 10}
	slotBody, _ := json.Marshal(slotPayload)
	slotReq, _ := http.NewRequest("POST", server.URL+"/api/admin/slots", bytes.NewReader(slotBody))
	slotReq.Header.Set("Content-Type", "application/json")
	slotReq.Header.Set("X-CSRF-Token", csrfToken)
	slotReq.Header.Set("Cookie", cookie)
	slotResp, err := http.DefaultClient.Do(slotReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, slotResp.StatusCode)

	var slotResult map[string]interface{}
	if err := json.NewDecoder(slotResp.Body).Decode(&slotResult); err != nil {
		t.Fatal(err)
	}
	slot := slotResult["slot"].(map[string]interface{})
	assert.Equal(t, "2025-12-15", slot["date"])
	assert.Equal(t, float64(10), slot["capacity"])

	// Test invalid slot creation (missing date)
	invalidSlotPayload := map[string]interface{}{"capacity": 10}
	invalidSlotBody, _ := json.Marshal(invalidSlotPayload)
	invalidSlotReq, _ := http.NewRequest("POST", server.URL+"/api/admin/slots", bytes.NewReader(invalidSlotBody))
	invalidSlotReq.Header.Set("Content-Type", "application/json")
	invalidSlotReq.Header.Set("X-CSRF-Token", csrfToken)
	invalidSlotReq.Header.Set("Cookie", cookie)
	invalidSlotResp, err := http.DefaultClient.Do(invalidSlotReq)
	assert.NoError(t, err)
	// Should fail due to missing date validation
	assert.Equal(t, 400, invalidSlotResp.StatusCode)
}

func TestAdminReservationManagement(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "adminuser3@example.com"
	password := "testpassword123"
	csrfToken, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Create a slot
	slotID := createSlot(server, csrfToken, cookie, "2025-12-25", 8)
	assert.Greater(t, slotID, 0)

	// Create a child
	childID := createChild(server, csrfToken, cookie, "TestChild", 5)
	assert.Greater(t, childID, 0)

	// Create a reservation
	resPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	resBody, _ := json.Marshal(resPayload)
	resReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	_, err := http.DefaultClient.Do(resReq)
	assert.NoError(t, err)

	// Get reservations list
	listReq, _ := http.NewRequest("GET", server.URL+"/api/admin/reservations", nil)
	listReq.Header.Set("Cookie", cookie)
	listResp, err := http.DefaultClient.Do(listReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, listResp.StatusCode)

	var listResult map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listResult); err != nil {
		t.Fatal(err)
	}

	reservations := listResult["data"].([]interface{})
	assert.Greater(t, len(reservations), 0)

	reservation := reservations[0].(map[string]interface{})
	reservationID := int(reservation["id"].(float64))

	// Test approve reservation
	approveReq, _ := http.NewRequest("PUT", server.URL+"/api/admin/approve/"+strconv.Itoa(reservationID), nil)
	approveReq.Header.Set("X-CSRF-Token", csrfToken)
	approveReq.Header.Set("Cookie", cookie)
	approveResp, err := http.DefaultClient.Do(approveReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, approveResp.StatusCode)

	var approveResult map[string]interface{}
	if err := json.NewDecoder(approveResp.Body).Decode(&approveResult); err != nil {
		t.Fatal(err)
	}

	reservationData := approveResult["reservation"].(map[string]interface{})
	assert.Equal(t, "approved", reservationData["status"])

	// Test reject reservation (use the same reservation since double booking is prevented)
	// Get the reservation again to make sure we have the latest data
	listReq2, _ := http.NewRequest("GET", server.URL+"/api/admin/reservations", nil)
	listReq2.Header.Set("Cookie", cookie)
	listResp2, err := http.DefaultClient.Do(listReq2)
	assert.NoError(t, err)
	var listResult2 map[string]interface{}
	if err := json.NewDecoder(listResp2.Body).Decode(&listResult2); err != nil {
		t.Fatal(err)
	}

	reservations2 := listResult2["data"].([]interface{})
	assert.Greater(t, len(reservations2), 0)
	reservation2 := reservations2[0].(map[string]interface{})
	reservationID2 := int(reservation2["id"].(float64))

	// Test reject reservation
	rejectReq, _ := http.NewRequest("PUT", server.URL+"/api/admin/reject/"+strconv.Itoa(reservationID2), nil)
	rejectReq.Header.Set("X-CSRF-Token", csrfToken)
	rejectReq.Header.Set("Cookie", cookie)
	rejectResp, err := http.DefaultClient.Do(rejectReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, rejectResp.StatusCode)

	var rejectResult map[string]interface{}
	if err := json.NewDecoder(rejectResp.Body).Decode(&rejectResult); err != nil {
		t.Fatal(err)
	}

	reservationData2 := rejectResult["reservation"].(map[string]interface{})
	assert.Equal(t, "rejected", reservationData2["status"])

	// Test approve/reject non-existent reservation
	approveNotFoundReq, _ := http.NewRequest("PUT", server.URL+"/api/admin/approve/99999", nil)
	approveNotFoundReq.Header.Set("X-CSRF-Token", csrfToken)
	approveNotFoundReq.Header.Set("Cookie", cookie)
	approveNotFoundResp, err := http.DefaultClient.Do(approveNotFoundReq)
	assert.NoError(t, err)
	assert.Equal(t, 404, approveNotFoundResp.StatusCode)

	rejectNotFoundReq, _ := http.NewRequest("PUT", server.URL+"/api/admin/reject/99999", nil)
	rejectNotFoundReq.Header.Set("X-CSRF-Token", csrfToken)
	rejectNotFoundReq.Header.Set("Cookie", cookie)
	rejectNotFoundResp, err := http.DefaultClient.Do(rejectNotFoundReq)
	assert.NoError(t, err)
	assert.Equal(t, 404, rejectNotFoundResp.StatusCode)
}

func TestAdminGetReservationsByStatus(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "adminuser4@example.com"
	password := "testpassword123"
	_, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Test get all reservations
	allReq, _ := http.NewRequest("GET", server.URL+"/api/admin/reservations", nil)
	allReq.Header.Set("Cookie", cookie)
	allResp, err := http.DefaultClient.Do(allReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, allResp.StatusCode)

	var allResult map[string]interface{}
	if err := json.NewDecoder(allResp.Body).Decode(&allResult); err != nil {
		t.Fatal(err)
	}

	// Debug: print the actual response structure
	t.Logf("Response structure: %+v", allResult)

	// Check if data field exists and is not nil
	dataField, exists := allResult["data"]
	if !exists {
		t.Fatalf("Response missing 'data' field. Full response: %+v", allResult)
	}
	if dataField == nil {
		t.Fatalf("Data field is nil. Full response: %+v", allResult)
	}

	allReservations := dataField.([]interface{})
	// Should be empty initially
	assert.Equal(t, 0, len(allReservations))

	// Test get reservations by status
	pendingReq, _ := http.NewRequest("GET", server.URL+"/api/admin/reservations?status=pending", nil)
	pendingReq.Header.Set("Cookie", cookie)
	pendingResp, err := http.DefaultClient.Do(pendingReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, pendingResp.StatusCode)

	var pendingResult map[string]interface{}
	if err := json.NewDecoder(pendingResp.Body).Decode(&pendingResult); err != nil {
		t.Fatal(err)
	}

	pendingReservations := pendingResult["data"].([]interface{})
	assert.Equal(t, 0, len(pendingReservations))

	// Test get approved reservations
	approvedReq, _ := http.NewRequest("GET", server.URL+"/api/admin/reservations?status=approved", nil)
	approvedReq.Header.Set("Cookie", cookie)
	approvedResp, err := http.DefaultClient.Do(approvedReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, approvedResp.StatusCode)

	var approvedResult map[string]interface{}
	if err := json.NewDecoder(approvedResp.Body).Decode(&approvedResult); err != nil {
		t.Fatal(err)
	}

	approvedReservations := approvedResult["data"].([]interface{})
	assert.Equal(t, 0, len(approvedReservations))
}

func TestAdminUserManagement(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "adminuser5@example.com"
	password := "testpassword123"
	csrfToken, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")

	// Test list users
	listReq, _ := http.NewRequest("GET", server.URL+"/api/admin/users", nil)
	listReq.Header.Set("Cookie", cookie)
	listResp, err := http.DefaultClient.Do(listReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, listResp.StatusCode)

	var listResult map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listResult); err != nil {
		t.Fatal(err)
	}

	users := listResult["data"].([]interface{})
	// Should have at least the admin user
	assert.GreaterOrEqual(t, len(users), 1)

	// Find a user to test with (not the current admin)
	var testUserID int
	for _, userInterface := range users {
		user := userInterface.(map[string]interface{})
		if user["email"] != email {
			testUserID = int(user["id"].(float64))
			break
		}
	}

	if testUserID == 0 {
		// Create a test user if none exists
		testUserPayload := map[string]string{"email": "testuser@example.com", "password": "testpassword123"}
		testUserBody, _ := json.Marshal(testUserPayload)
		testUserReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(testUserBody))
		testUserReq.Header.Set("Content-Type", "application/json")
		testUserReq.Header.Set("X-CSRF-Token", csrfToken)
		testUserReq.Header.Set("Cookie", cookie)
		_, err := http.DefaultClient.Do(testUserReq)
		assert.NoError(t, err)

		// Get the user ID
		listReq2, _ := http.NewRequest("GET", server.URL+"/api/admin/users", nil)
		listReq2.Header.Set("Cookie", cookie)
		listResp2, err := http.DefaultClient.Do(listReq2)
		assert.NoError(t, err)
		var listResult2 map[string]interface{}
		if err := json.NewDecoder(listResp2.Body).Decode(&listResult2); err != nil {
			t.Fatal(err)
		}

		users2 := listResult2["data"].([]interface{})
		for _, userInterface := range users2 {
			user := userInterface.(map[string]interface{})
			if user["email"] == "testuser@example.com" {
				testUserID = int(user["id"].(float64))
				break
			}
		}
	}

	// Test update user role
	rolePayload := map[string]interface{}{"role": "admin"}
	roleBody, _ := json.Marshal(rolePayload)
	roleReq, _ := http.NewRequest("PUT", server.URL+"/api/admin/users/"+strconv.Itoa(testUserID)+"/role", bytes.NewReader(roleBody))
	roleReq.Header.Set("Content-Type", "application/json")
	roleReq.Header.Set("X-CSRF-Token", csrfToken)
	roleReq.Header.Set("Cookie", cookie)
	roleResp, err := http.DefaultClient.Do(roleReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, roleResp.StatusCode)

	var roleResult map[string]interface{}
	if err := json.NewDecoder(roleResp.Body).Decode(&roleResult); err != nil {
		t.Fatal(err)
	}

	userData := roleResult["user"].(map[string]interface{})
	assert.Equal(t, "admin", userData["role"])

	// Test delete user
	deleteReq, _ := http.NewRequest("DELETE", server.URL+"/api/admin/users/"+strconv.Itoa(testUserID), nil)
	deleteReq.Header.Set("X-CSRF-Token", csrfToken)
	deleteReq.Header.Set("Cookie", cookie)
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, deleteResp.StatusCode)

	var deleteResult map[string]interface{}
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleteResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "User deleted successfully", deleteResult["message"])

	// Test delete non-existent user
	deleteNotFoundReq, _ := http.NewRequest("DELETE", server.URL+"/api/admin/users/99999", nil)
	deleteNotFoundReq.Header.Set("X-CSRF-Token", csrfToken)
	deleteNotFoundReq.Header.Set("Cookie", cookie)
	deleteNotFoundResp, err := http.DefaultClient.Do(deleteNotFoundReq)
	assert.NoError(t, err)
	assert.Equal(t, 404, deleteNotFoundResp.StatusCode) // Should return 404 for non-existent user
}

func TestAdminUnauthorizedAccess(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	email := "parentuser2@example.com"
	password := "testpassword123"
	csrfToken, cookie := registerAndLogin(server, email, password, initToken, initCookie)

	// Ensure user is parent (not admin)
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Test admin endpoints with parent role (should be unauthorized)
	// Use appropriate HTTP methods for each endpoint
	testCases := []struct {
		endpoint string
		method   string
	}{
		{"/api/admin/slots", "POST"},
		{"/api/admin/users", "GET"},
		{"/api/admin/reservations", "GET"},
	}

	for _, tc := range testCases {
		var req *http.Request
		if tc.method == "POST" {
			// For POST requests, send some data
			payload := map[string]interface{}{"date": "2025-12-25", "capacity": 5}
			body, _ := json.Marshal(payload)
			req, _ = http.NewRequest(tc.method, server.URL+tc.endpoint, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-CSRF-Token", csrfToken)
		} else {
			req, _ = http.NewRequest(tc.method, server.URL+tc.endpoint, nil)
		}
		req.Header.Set("Cookie", cookie)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		// Application returns 403 for admin endpoints when user is not admin
		assert.Equal(t, 403, resp.StatusCode, "Expected 403 for endpoint: %s with method %s", tc.endpoint, tc.method)
	}
}
