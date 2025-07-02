package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	payload := map[string]string{
		"email":    "testuser@example.com",
		"password": "testpassword123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "User registered", result["message"])
}

func TestLogin_Success(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	payload := map[string]string{"email": "loginuser@example.com", "password": "testpassword123"}
	body, _ := json.Marshal(payload)
	// Register first
	regReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		regReq.Header.Set("Cookie", cookie)
	}
	app.Test(regReq, -1)
	// Now login
	loginBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Logged in", result["message"])
}

func TestLogin_Failure(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	payload := map[string]string{"email": "nouser@example.com", "password": "wrongpass"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Invalid credentials", result["error"])
}

func TestRegister_InvalidEmail(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	payload := map[string]string{"email": "notanemail", "password": "testpassword123"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Invalid email format", result["error"])
}

func TestGetProfile(t *testing.T) {
	app := setupTestApp()
	// Register and login user (CSRF not needed for GET)
	payload := map[string]string{"email": "profileuser@example.com", "password": "testpassword123"}
	body, _ := json.Marshal(payload)
	app.Test(httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body)), -1)
	// Simulate session (mocking session is more complex, so this is a placeholder for now)
	// In a real test, you would mock session or use Fiber's session middleware with a test store
}

func TestUserProfileEndpoints(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "profileuser@example.com", "testpassword123", csrfToken, cookie)

	// Get profile
	getReq := httptest.NewRequest("GET", "/api/user/profile", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := app.Test(getReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, getResp.StatusCode)
	var profile map[string]interface{}
	json.NewDecoder(getResp.Body).Decode(&profile)
	assert.Equal(t, "profileuser@example.com", profile["Email"])

	// Update profile
	updatePayload := map[string]interface{}{"email": "profileuser2@example.com", "password": "newpassword123"}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest("PUT", "/api/user/profile", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("X-CSRF-Token", csrfToken)
	updateReq.Header.Set("Cookie", cookie)
	updateResp, err := app.Test(updateReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, updateResp.StatusCode)
	var updateResult map[string]interface{}
	json.NewDecoder(updateResp.Body).Decode(&updateResult)
	assert.Equal(t, "Profile updated", updateResult["message"])
}

func TestPasswordResetEndpoints(t *testing.T) {
	app := setupTestApp()
	csrfToken, cookie := getCSRFTokenAndCookie(app)
	registerAndLogin(app, "resetuser@example.com", "testpassword123", csrfToken, cookie)

	// Request password reset
	resetPayload := map[string]interface{}{"email": "resetuser@example.com"}
	resetBody, _ := json.Marshal(resetPayload)
	resetReq := httptest.NewRequest("POST", "/api/auth/request-reset", bytes.NewReader(resetBody))
	resetReq.Header.Set("Content-Type", "application/json")
	resetReq.Header.Set("X-CSRF-Token", csrfToken)
	resetReq.Header.Set("Cookie", cookie)
	resetResp, err := app.Test(resetReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resetResp.StatusCode)
	var resetResult map[string]interface{}
	json.NewDecoder(resetResp.Body).Decode(&resetResult)
	assert.Contains(t, resetResult["message"], "Password reset email sent")

	// Simulate token (would be sent by email in real app)
	// For test, get token from controllers package (if exported) or skip actual reset
}

func TestHealthAndVersionEndpoints(t *testing.T) {
	app := setupTestApp()

	// Health
	healthReq := httptest.NewRequest("GET", "/health", nil)
	healthResp, err := app.Test(healthReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, healthResp.StatusCode)
	buf := new(bytes.Buffer)
	buf.ReadFrom(healthResp.Body)
	assert.Equal(t, "OK", buf.String())

	// Version
	versionReq := httptest.NewRequest("GET", "/version", nil)
	versionResp, err := app.Test(versionReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, versionResp.StatusCode)
	var version map[string]interface{}
	json.NewDecoder(versionResp.Body).Decode(&version)
	assert.Equal(t, "1.0.0", version["version"])
}
