package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reservio/config"
	"reservio/models"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func promoteToAdmin(email string) {
	var user models.User
	config.DB.Where("email = ?", email).First(&user)
	user.Role = "admin"
	config.DB.Save(&user)
}

func TestAdminSlotCreationAndUserManagement(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "adminuser@example.com", "testpassword123", csrfToken, cookie)
	promoteToAdmin("adminuser@example.com")

	// Create slot
	slotPayload := map[string]interface{}{"date": "2025-12-02", "capacity": 20}
	slotBody, _ := json.Marshal(slotPayload)
	slotReq := httptest.NewRequest("POST", "/api/admin/slots", bytes.NewReader(slotBody))
	slotReq.Header.Set("Content-Type", "application/json")
	slotReq.Header.Set("X-CSRF-Token", csrfToken)
	slotReq.Header.Set("Cookie", cookie)
	slotResp, err := app.Test(slotReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, slotResp.StatusCode)
	var slotResult map[string]interface{}
	json.NewDecoder(slotResp.Body).Decode(&slotResult)
	assert.Equal(t, "2025-12-02", slotResult["Date"])

	// List users
	usersReq := httptest.NewRequest("GET", "/api/admin/users", nil)
	usersReq.Header.Set("Cookie", cookie)
	usersResp, err := app.Test(usersReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, usersResp.StatusCode)
	var users []map[string]interface{}
	json.NewDecoder(usersResp.Body).Decode(&users)
	assert.True(t, len(users) > 0)

	// Change user role
	uid := int(users[0]["ID"].(float64))
	rolePayload := map[string]interface{}{"role": "admin"}
	roleBody, _ := json.Marshal(rolePayload)
	roleReq := httptest.NewRequest("PUT", "/api/admin/users/"+strconv.Itoa(uid)+"/role", bytes.NewReader(roleBody))
	roleReq.Header.Set("Content-Type", "application/json")
	roleReq.Header.Set("X-CSRF-Token", csrfToken)
	roleReq.Header.Set("Cookie", cookie)
	roleResp, err := app.Test(roleReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, roleResp.StatusCode)
}

func TestAdminOnlyAccess(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "parentuser@example.com", "testpassword123", csrfToken, cookie)
	// Try to create slot as non-admin
	slotPayload := map[string]interface{}{"date": "2025-12-03", "capacity": 10}
	slotBody, _ := json.Marshal(slotPayload)
	slotReq := httptest.NewRequest("POST", "/api/admin/slots", bytes.NewReader(slotBody))
	slotReq.Header.Set("Content-Type", "application/json")
	slotReq.Header.Set("X-CSRF-Token", csrfToken)
	slotReq.Header.Set("Cookie", cookie)
	slotResp, err := app.Test(slotReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 403, slotResp.StatusCode)
}
