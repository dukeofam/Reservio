package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reservio/config"
	"strconv"
	"testing"

	"reservio/models"

	"github.com/stretchr/testify/assert"
)

func TestMakeListCancelReservation(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "resparent+1@example.com", "testpassword123", csrfToken, cookie)

	// Add child
	childPayload := map[string]interface{}{"name": "Bob", "age": 4}
	childBody, _ := json.Marshal(childPayload)
	addReq := httptest.NewRequest("POST", "/api/parent/children", bytes.NewReader(childBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, _ := app.Test(addReq, -1)
	var addedChild map[string]interface{}
	json.NewDecoder(addResp.Body).Decode(&addedChild)
	childID := int(addedChild["ID"].(float64))

	// Insert slot directly into DB
	slot := models.Slot{Date: "2025-12-01", Capacity: 10}
	config.DB.Create(&slot)

	// Get slot ID
	slotsReq := httptest.NewRequest("GET", "/api/slots", nil)
	slotsReq.Header.Set("Cookie", cookie)
	slotsResp, _ := app.Test(slotsReq, -1)
	var slots []map[string]interface{}
	json.NewDecoder(slotsResp.Body).Decode(&slots)
	if len(slots) == 0 {
		t.Fatalf("No slots found. Make sure a slot is created before making a reservation.")
	}
	slotID := int(slots[0]["ID"].(float64))

	// Make reservation
	resPayload := map[string]interface{}{"child_id": childID, "slot_id": slotID}
	resBody, _ := json.Marshal(resPayload)
	resReq := httptest.NewRequest("POST", "/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	resResp, err := app.Test(resReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resResp.StatusCode)
	var resResult map[string]interface{}
	json.NewDecoder(resResp.Body).Decode(&resResult)
	assert.Equal(t, "Reservation requested", resResult["message"])

	// List reservations
	listReq := httptest.NewRequest("GET", "/api/parent/reservations", nil)
	listReq.Header.Set("Cookie", cookie)
	listResp, err := app.Test(listReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, listResp.StatusCode)
	var reservations []map[string]interface{}
	json.NewDecoder(listResp.Body).Decode(&reservations)
	assert.True(t, len(reservations) > 0)
	reservationID := int(reservations[0]["ID"].(float64))

	// Cancel reservation
	cancelReq := httptest.NewRequest("DELETE", "/api/parent/reservations/"+strconv.Itoa(reservationID), nil)
	cancelReq.Header.Set("X-CSRF-Token", csrfToken)
	cancelReq.Header.Set("Cookie", cookie)
	cancelResp, err := app.Test(cancelReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, cancelResp.StatusCode)
	var cancelResult map[string]interface{}
	json.NewDecoder(cancelResp.Body).Decode(&cancelResult)
	assert.Equal(t, "Reservation cancelled", cancelResult["message"])
}

func TestMakeReservation_InvalidData(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "resparent+2@example.com", "testpassword123", csrfToken, cookie)
	// Missing child_id
	resPayload := map[string]interface{}{"slot_id": 1}
	resBody, _ := json.Marshal(resPayload)
	resReq := httptest.NewRequest("POST", "/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	resResp, err := app.Test(resReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, resResp.StatusCode)
}

func TestAdminReservationApprovalRejection(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "admin2@example.com", "testpassword123", csrfToken, cookie)
	promoteToAdmin("admin2@example.com")

	// Create slot
	slotID := createSlot(app, csrfToken, cookie, "2025-12-10", 5)

	// Register and login as parent, make reservation
	csrf2, cookie2 := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "parent2@example.com", "testpassword123", csrf2, cookie2)
	childID := createChild(app, csrf2, cookie2, "TestChild", "2018-01-01")
	resID := createReservation(app, csrf2, cookie2, slotID, childID)

	// Approve reservation as admin (PUT /api/admin/approve/:id)
	approveReq := httptest.NewRequest("PUT", "/api/admin/approve/"+strconv.Itoa(resID), nil)
	approveReq.Header.Set("X-CSRF-Token", csrfToken)
	approveReq.Header.Set("Cookie", cookie)
	approveResp, err := app.Test(approveReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, approveResp.StatusCode)

	// Reject reservation as admin (PUT /api/admin/reject/:id)
	resID2 := createReservation(app, csrf2, cookie2, slotID, childID)
	rejectReq := httptest.NewRequest("PUT", "/api/admin/reject/"+strconv.Itoa(resID2), nil)
	rejectReq.Header.Set("X-CSRF-Token", csrfToken)
	rejectReq.Header.Set("Cookie", cookie)
	rejectResp, err := app.Test(rejectReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, rejectResp.StatusCode)
}

func TestSlotListingAndEdgeCases(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "parent3@example.com", "testpassword123", csrfToken, cookie)
	// Create slot as admin
	promoteToAdmin("parent3@example.com")
	createSlot(app, csrfToken, cookie, "2025-12-20", 10)

	// List slots as parent
	slotsReq := httptest.NewRequest("GET", "/api/slots", nil)
	slotsReq.Header.Set("Cookie", cookie)
	slotsResp, err := app.Test(slotsReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, slotsResp.StatusCode)
	var slots []map[string]interface{}
	json.NewDecoder(slotsResp.Body).Decode(&slots)
	assert.True(t, len(slots) > 0)

	// List slots as unauthenticated user
	slotsReq2 := httptest.NewRequest("GET", "/api/slots", nil)
	slotsResp2, err := app.Test(slotsReq2, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, slotsResp2.StatusCode)
}

func TestUnauthorizedAndInvalidCSRF(t *testing.T) {
	app := setupTestApp()
	// No cookie, no CSRF
	req := httptest.NewRequest("POST", "/api/children", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	// Valid cookie, missing CSRF
	csrf, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "parent4@example.com", "testpassword123", csrf, cookie)
	childPayload := map[string]interface{}{"name": "NoCSRF", "birthdate": "2017-01-01"}
	childBody, _ := json.Marshal(childPayload)
	childReq := httptest.NewRequest("POST", "/api/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("Cookie", cookie)
	childResp, err := app.Test(childReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 403, childResp.StatusCode)

	// Invalid session
	childReq2 := httptest.NewRequest("POST", "/api/children", bytes.NewReader(childBody))
	childReq2.Header.Set("Content-Type", "application/json")
	childReq2.Header.Set("X-CSRF-Token", csrf)
	childReq2.Header.Set("Cookie", "invalidsession=123")
	childResp2, err := app.Test(childReq2, -1)
	assert.NoError(t, err)
	assert.Equal(t, 403, childResp2.StatusCode)
}
