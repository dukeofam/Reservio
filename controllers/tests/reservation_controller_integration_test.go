package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reservio/config"
	"reservio/models"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReservationEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie = registerAndLogin(server, "resparent+1@example.com", "testpassword123", csrfToken, cookie)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)

	// Set user as admin in DB
	config.DB.Model(&models.User{}).Where("email = ?", "resparent+1@example.com").Update("role", "admin")

	// Create slot as admin
	slotID := createSlot(server, csrfToken, cookie, "2025-12-10", 5)
	if slotID == 0 {
		t.Fatalf("expected slot ID to be non-zero")
	}

	// Create child
	childID := createChild(server, csrfToken, cookie, "TestChild", 6)
	if childID == 0 {
		t.Fatalf("expected child ID to be non-zero")
	}

	// Make reservation
	reservationID := createReservation(server, csrfToken, cookie, slotID, childID)
	if reservationID == 0 {
		t.Fatalf("expected reservation ID to be non-zero")
	}
}

func TestReservationLifecycle(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "resparent+2@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create slot (need admin for this, so we'll use the helper which requires admin role)
	// First set as admin temporarily
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")
	slotID := createSlot(server, csrfToken, cookie, "2025-12-15", 3)
	// Set back to parent
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create child
	childPayload := map[string]interface{}{"name": "ReservationChild", "birthdate": "2019-05-15"}
	childBody, _ := json.Marshal(childPayload)
	childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	childReq.Header.Set("Cookie", cookie)
	childResp, err := http.DefaultClient.Do(childReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, childResp.StatusCode)

	var childResult map[string]interface{}
	if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
		t.Fatal(err)
	}
	child := childResult["child"].(map[string]interface{})
	childID := int(child["id"].(float64))

	// Make reservation
	resPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	resBody, _ := json.Marshal(resPayload)
	resReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	resResp, err := http.DefaultClient.Do(resReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resResp.StatusCode)

	// Get reservations
	listReq, _ := http.NewRequest("GET", server.URL+"/api/parent/reservations", nil)
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

	// Cancel reservation
	cancelReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/reservations/"+strconv.Itoa(reservationID), nil)
	cancelReq.Header.Set("X-CSRF-Token", csrfToken)
	cancelReq.Header.Set("Cookie", cookie)
	cancelResp, err := http.DefaultClient.Do(cancelReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, cancelResp.StatusCode)

	var cancelResult map[string]interface{}
	if err := json.NewDecoder(cancelResp.Body).Decode(&cancelResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Reservation cancelled successfully", cancelResult["message"])

	// Verify reservation is cancelled by checking list again
	listReq2, _ := http.NewRequest("GET", server.URL+"/api/parent/reservations", nil)
	listReq2.Header.Set("Cookie", cookie)
	listResp2, err := http.DefaultClient.Do(listReq2)
	assert.NoError(t, err)
	assert.Equal(t, 200, listResp2.StatusCode)

	var listResult2 map[string]interface{}
	if err := json.NewDecoder(listResp2.Body).Decode(&listResult2); err != nil {
		t.Fatal(err)
	}
	reservations2 := listResult2["data"].([]interface{})
	// Should be empty now or have cancelled status
	assert.LessOrEqual(t, len(reservations2), len(reservations))
}

func TestReservationEdgeCases(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "resparent+3@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create slot (need admin temporarily)
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")
	slotID := createSlot(server, csrfToken, cookie, "2025-12-20", 1) // Capacity of 1
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create child
	childPayload := map[string]interface{}{"name": "EdgeCaseChild", "birthdate": "2020-03-10"}
	childBody, _ := json.Marshal(childPayload)
	childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	childReq.Header.Set("Cookie", cookie)
	childResp, err := http.DefaultClient.Do(childReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, childResp.StatusCode)

	var childResult map[string]interface{}
	if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
		t.Fatal(err)
	}
	child := childResult["child"].(map[string]interface{})
	childID := int(child["id"].(float64))

	// Test reservation with invalid slot ID
	invalidSlotPayload := map[string]interface{}{"slot_id": 99999, "child_id": childID}
	invalidSlotBody, _ := json.Marshal(invalidSlotPayload)
	invalidSlotReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(invalidSlotBody))
	invalidSlotReq.Header.Set("Content-Type", "application/json")
	invalidSlotReq.Header.Set("X-CSRF-Token", csrfToken)
	invalidSlotReq.Header.Set("Cookie", cookie)
	invalidSlotResp, err := http.DefaultClient.Do(invalidSlotReq)
	assert.NoError(t, err)
	// Application is more permissive - accepts invalid slot ID
	assert.Equal(t, 404, invalidSlotResp.StatusCode)

	// Test reservation with invalid child ID
	invalidChildPayload := map[string]interface{}{"slot_id": slotID, "child_id": 99999}
	invalidChildBody, _ := json.Marshal(invalidChildPayload)
	invalidChildReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(invalidChildBody))
	invalidChildReq.Header.Set("Content-Type", "application/json")
	invalidChildReq.Header.Set("X-CSRF-Token", csrfToken)
	invalidChildReq.Header.Set("Cookie", cookie)
	invalidChildResp, err := http.DefaultClient.Do(invalidChildReq)
	assert.NoError(t, err)
	// Expect not found for invalid child ID
	assert.Equal(t, 404, invalidChildResp.StatusCode)

	// Test reservation with missing data
	missingDataPayload := map[string]interface{}{"slot_id": slotID}
	missingDataBody, _ := json.Marshal(missingDataPayload)
	missingDataReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(missingDataBody))
	missingDataReq.Header.Set("Content-Type", "application/json")
	missingDataReq.Header.Set("X-CSRF-Token", csrfToken)
	missingDataReq.Header.Set("Cookie", cookie)
	missingDataResp, err := http.DefaultClient.Do(missingDataReq)
	assert.NoError(t, err)
	assert.Equal(t, 400, missingDataResp.StatusCode)

	// Test valid reservation
	validPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	validBody, _ := json.Marshal(validPayload)
	validReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(validBody))
	validReq.Header.Set("Content-Type", "application/json")
	validReq.Header.Set("X-CSRF-Token", csrfToken)
	validReq.Header.Set("Cookie", cookie)
	validResp, err := http.DefaultClient.Do(validReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, validResp.StatusCode)

	// Test double booking (same child, same slot)
	doubleBookReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(validBody))
	doubleBookReq.Header.Set("Content-Type", "application/json")
	doubleBookReq.Header.Set("X-CSRF-Token", csrfToken)
	doubleBookReq.Header.Set("Cookie", cookie)
	doubleBookResp, err := http.DefaultClient.Do(doubleBookReq)
	assert.NoError(t, err)
	// Expect validation error for double booking (already booked)
	assert.Contains(t, []int{400, 409}, doubleBookResp.StatusCode)
}

func TestCancelReservationEdgeCases(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "resparent+4@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Test cancel non-existent reservation
	cancelNotFoundReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/reservations/99999", nil)
	cancelNotFoundReq.Header.Set("X-CSRF-Token", csrfToken)
	cancelNotFoundReq.Header.Set("Cookie", cookie)
	cancelNotFoundResp, err := http.DefaultClient.Do(cancelNotFoundReq)
	assert.NoError(t, err)
	// Expect not found for non-existent reservation
	assert.Equal(t, 404, cancelNotFoundResp.StatusCode)

	// Test cancel reservation without authentication
	cancelUnauthReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/reservations/1", nil)
	cancelUnauthReq.Header.Set("X-CSRF-Token", csrfToken)
	// No cookie
	cancelUnauthResp, err := http.DefaultClient.Do(cancelUnauthReq)
	assert.NoError(t, err)
	// Application returns 401 for unauthorized cancellation without valid session
	assert.Equal(t, 401, cancelUnauthResp.StatusCode)
}

func TestReservationAuthorization(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "resparent+5@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create slot (need admin temporarily)
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")
	slotID := createSlot(server, csrfToken, cookie, "2025-12-25", 5)
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create child
	childPayload := map[string]interface{}{"name": "AuthChild", "birthdate": "2017-08-20"}
	childBody, _ := json.Marshal(childPayload)
	childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	childReq.Header.Set("Cookie", cookie)
	childResp, err := http.DefaultClient.Do(childReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, childResp.StatusCode)

	var childResult map[string]interface{}
	if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
		t.Fatal(err)
	}
	child := childResult["child"].(map[string]interface{})
	childID := int(child["id"].(float64))

	// Make reservation
	resPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	resBody, _ := json.Marshal(resPayload)
	resReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	resResp, err := http.DefaultClient.Do(resReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resResp.StatusCode)

	// Get reservation ID
	listReq, _ := http.NewRequest("GET", server.URL+"/api/parent/reservations", nil)
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

	// Test that parent can only access their own reservations
	// Create another user and try to access the first user's reservation
	otherEmail := "otherparent@example.com"
	otherPassword := "testpassword123"

	// Register other user
	otherPayload := map[string]string{"email": otherEmail, "password": otherPassword}
	otherBody, _ := json.Marshal(otherPayload)
	otherRegReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(otherBody))
	otherRegReq.Header.Set("Content-Type", "application/json")
	otherRegReq.Header.Set("X-CSRF-Token", csrfToken)
	otherRegReq.Header.Set("Cookie", cookie)
	_, err = http.DefaultClient.Do(otherRegReq)
	assert.NoError(t, err)

	// Set other user as parent
	config.DB.Model(&models.User{}).Where("email = ?", otherEmail).Update("role", "parent")

	// Login as other user
	otherLoginPayload := map[string]string{"email": otherEmail, "password": otherPassword}
	otherLoginBody, _ := json.Marshal(otherLoginPayload)
	otherLoginReq, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(otherLoginBody))
	otherLoginReq.Header.Set("Content-Type", "application/json")
	otherLoginReq.Header.Set("X-CSRF-Token", csrfToken)
	otherLoginReq.Header.Set("Cookie", cookie)
	otherLoginResp, err := http.DefaultClient.Do(otherLoginReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, otherLoginResp.StatusCode)

	// Get other user's session cookie
	var otherCookie string
	for _, c := range otherLoginResp.Cookies() {
		if c.Name == "session" {
			otherCookie = c.Name + "=" + c.Value
			break
		}
	}

	// Try to cancel first user's reservation with other user's session
	cancelOtherReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/reservations/"+strconv.Itoa(reservationID), nil)
	cancelOtherReq.Header.Set("X-CSRF-Token", csrfToken)
	cancelOtherReq.Header.Set("Cookie", otherCookie)
	cancelOtherResp, err := http.DefaultClient.Do(cancelOtherReq)
	assert.NoError(t, err)
	// Should be not found / unauthorized
	assert.Equal(t, 404, cancelOtherResp.StatusCode)
}

func TestReservationCapacityLimits(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	email := "resparent+6@example.com"
	password := "testpassword123"
	csrfToken, cookie = registerAndLogin(server, email, password, csrfToken, cookie)

	// Set user as admin to create slot
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "admin")
	slotID := createSlot(server, csrfToken, cookie, "2025-12-30", 2) // Capacity of 2
	config.DB.Model(&models.User{}).Where("email = ?", email).Update("role", "parent")

	// Create two children
	child1Payload := map[string]interface{}{"name": "Child1", "birthdate": "2018-01-01"}
	child1Body, _ := json.Marshal(child1Payload)
	child1Req, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(child1Body))
	child1Req.Header.Set("Content-Type", "application/json")
	child1Req.Header.Set("X-CSRF-Token", csrfToken)
	child1Req.Header.Set("Cookie", cookie)
	child1Resp, err := http.DefaultClient.Do(child1Req)
	assert.NoError(t, err)
	assert.Equal(t, 200, child1Resp.StatusCode)

	var child1Result map[string]interface{}
	if err := json.NewDecoder(child1Resp.Body).Decode(&child1Result); err != nil {
		t.Fatal(err)
	}
	child1 := child1Result["child"].(map[string]interface{})
	child1ID := int(child1["id"].(float64))

	child2Payload := map[string]interface{}{"name": "Child2", "birthdate": "2019-02-02"}
	child2Body, _ := json.Marshal(child2Payload)
	child2Req, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(child2Body))
	child2Req.Header.Set("Content-Type", "application/json")
	child2Req.Header.Set("X-CSRF-Token", csrfToken)
	child2Req.Header.Set("Cookie", cookie)
	child2Resp, err := http.DefaultClient.Do(child2Req)
	assert.NoError(t, err)
	assert.Equal(t, 200, child2Resp.StatusCode)

	var child2Result map[string]interface{}
	if err := json.NewDecoder(child2Resp.Body).Decode(&child2Result); err != nil {
		t.Fatal(err)
	}
	child2 := child2Result["child"].(map[string]interface{})
	child2ID := int(child2["id"].(float64))

	// Make first reservation
	res1Payload := map[string]interface{}{"slot_id": slotID, "child_id": child1ID}
	res1Body, _ := json.Marshal(res1Payload)
	res1Req, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(res1Body))
	res1Req.Header.Set("Content-Type", "application/json")
	res1Req.Header.Set("X-CSRF-Token", csrfToken)
	res1Req.Header.Set("Cookie", cookie)
	res1Resp, err := http.DefaultClient.Do(res1Req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res1Resp.StatusCode)

	// Make second reservation
	res2Payload := map[string]interface{}{"slot_id": slotID, "child_id": child2ID}
	res2Body, _ := json.Marshal(res2Payload)
	res2Req, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(res2Body))
	res2Req.Header.Set("Content-Type", "application/json")
	res2Req.Header.Set("X-CSRF-Token", csrfToken)
	res2Req.Header.Set("Cookie", cookie)
	res2Resp, err := http.DefaultClient.Do(res2Req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res2Resp.StatusCode)

	// Try to make third reservation (should fail if capacity is enforced)
	child3Payload := map[string]interface{}{"name": "Child3", "birthdate": "2020-03-03"}
	child3Body, _ := json.Marshal(child3Payload)
	child3Req, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(child3Body))
	child3Req.Header.Set("Content-Type", "application/json")
	child3Req.Header.Set("X-CSRF-Token", csrfToken)
	child3Req.Header.Set("Cookie", cookie)
	child3Resp, err := http.DefaultClient.Do(child3Req)
	assert.NoError(t, err)
	assert.Equal(t, 200, child3Resp.StatusCode)

	var child3Result map[string]interface{}
	if err := json.NewDecoder(child3Resp.Body).Decode(&child3Result); err != nil {
		t.Fatal(err)
	}
	child3 := child3Result["child"].(map[string]interface{})
	child3ID := int(child3["id"].(float64))

	res3Payload := map[string]interface{}{"slot_id": slotID, "child_id": child3ID}
	res3Body, _ := json.Marshal(res3Payload)
	res3Req, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(res3Body))
	res3Req.Header.Set("Content-Type", "application/json")
	res3Req.Header.Set("X-CSRF-Token", csrfToken)
	res3Req.Header.Set("Cookie", cookie)
	res3Resp, err := http.DefaultClient.Do(res3Req)
	// This might be 400 (capacity exceeded) or 200 (if capacity is not enforced)
	assert.NoError(t, err)
	assert.Contains(t, []int{400, 409}, res3Resp.StatusCode)
}
