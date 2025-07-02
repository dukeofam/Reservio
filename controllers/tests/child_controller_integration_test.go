package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGetEditDeleteChild(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "childparent+1@example.com", "testpassword123", csrfToken, cookie)

	// Add child
	childPayload := map[string]interface{}{"name": "Alice", "age": 5}
	childBody, _ := json.Marshal(childPayload)
	addReq := httptest.NewRequest("POST", "/api/parent/children", bytes.NewReader(childBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, err := app.Test(addReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, addResp.StatusCode)
	var added map[string]interface{}
	if err := json.NewDecoder(addResp.Body).Decode(&added); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Alice", added["Name"])

	// Get children
	getReq := httptest.NewRequest("GET", "/api/parent/children", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := app.Test(getReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, getResp.StatusCode)
	var children []map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&children); err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(children) > 0)

	// Edit child
	childID := int(added["ID"].(float64))
	editPayload := map[string]interface{}{"name": "AliceUpdated", "age": 6}
	editBody, _ := json.Marshal(editPayload)
	editReq := httptest.NewRequest("PUT", "/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(editBody))
	editReq.Header.Set("Content-Type", "application/json")
	editReq.Header.Set("X-CSRF-Token", csrfToken)
	editReq.Header.Set("Cookie", cookie)
	editResp, err := app.Test(editReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, editResp.StatusCode)
	var edited map[string]interface{}
	if err := json.NewDecoder(editResp.Body).Decode(&edited); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "AliceUpdated", edited["Name"])

	// Delete child
	delReq := httptest.NewRequest("DELETE", "/api/parent/children/"+strconv.Itoa(childID), nil)
	delReq.Header.Set("X-CSRF-Token", csrfToken)
	delReq.Header.Set("Cookie", cookie)
	delResp, err := app.Test(delReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, delResp.StatusCode)
	var delResult map[string]interface{}
	if err := json.NewDecoder(delResp.Body).Decode(&delResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Child deleted", delResult["message"])
}

func TestAddChild_InvalidData(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "childparent+2@example.com", "testpassword123", csrfToken, cookie)
	// Missing name
	childPayload := map[string]interface{}{"age": 5}
	childBody, _ := json.Marshal(childPayload)
	addReq := httptest.NewRequest("POST", "/api/parent/children", bytes.NewReader(childBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, err := app.Test(addReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, addResp.StatusCode)
}

func TestDeleteChild_NotFound(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "childparent+3@example.com", "testpassword123", csrfToken, cookie)
	delReq := httptest.NewRequest("DELETE", "/api/parent/children/99999", nil)
	delReq.Header.Set("X-CSRF-Token", csrfToken)
	delReq.Header.Set("Cookie", cookie)
	delResp, err := app.Test(delReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 404, delResp.StatusCode)
}
