package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	payload := map[string]string{
		"email":    "testuser@example.com",
		"password": "testpassword123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["message"] != "User registered" {
		t.Fatalf("expected 'User registered', got %v", result["message"])
	}
}

func TestLogin_Success(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	payload := map[string]string{"email": "loginuser@example.com", "password": "testpassword123"}
	body, _ := json.Marshal(payload)
	// Register first
	regReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		regReq.Header.Set("Cookie", cookie)
	}
	if _, err := http.DefaultClient.Do(regReq); err != nil {
		t.Fatal(err)
	}
	// Now login
	loginBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Logged in", result["message"])
}

func TestLogin_Failure(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	payload := map[string]string{"email": "nouser@example.com", "password": "wrongpass"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Invalid credentials", result["error"])
}

func TestRegister_InvalidEmail(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	payload := map[string]string{"email": "notanemail", "password": "testpassword123"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Invalid email format", result["error"])
}

func TestGetProfile(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	_, cookie = registerAndLogin(server, "profileuser@example.com", "testpassword123", csrfToken, cookie)

	getReq, _ := http.NewRequest("GET", server.URL+"/api/user/profile", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("http.DefaultClient.Do error (get profile): %v", err)
	}
	if getResp.StatusCode != 200 {
		var bodyBytes bytes.Buffer
		_, _ = bodyBytes.ReadFrom(getResp.Body)
		t.Fatalf("expected 200, got %d, body: %s", getResp.StatusCode, bodyBytes.String())
	}
	var profile map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&profile); err != nil {
		t.Fatal(err)
	}
	if profile["email"] != "profileuser@example.com" {
		t.Fatalf("expected email 'profileuser@example.com', got %v", profile["email"])
	}
}

func TestUserProfileEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	csrfToken, cookie = registerAndLogin(server, "profileuser@example.com", "testpassword123", csrfToken, cookie)

	// Get profile
	getReq, _ := http.NewRequest("GET", server.URL+"/api/user/profile", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, getErr := http.DefaultClient.Do(getReq)
	if getErr != nil {
		t.Fatalf("http.DefaultClient.Do error (get profile): %v", getErr)
	}
	if getResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}
	var profile map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&profile); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "profileuser@example.com", profile["email"])

	// Update profile
	updatePayload := map[string]interface{}{"email": "profileuser2@example.com", "password": "newpassword123"}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq, _ := http.NewRequest("PUT", server.URL+"/api/user/profile", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("X-CSRF-Token", csrfToken)
	updateReq.Header.Set("Cookie", cookie)
	updateResp, updateErr := http.DefaultClient.Do(updateReq)
	if updateErr != nil {
		t.Fatalf("http.DefaultClient.Do error (update profile): %v", updateErr)
	}
	if updateResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", updateResp.StatusCode)
	}
	var updateResult map[string]interface{}
	if err := json.NewDecoder(updateResp.Body).Decode(&updateResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Profile updated", updateResult["message"])
}

func TestPasswordResetEndpoints(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping password reset test in CI")
	}
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	registerAndLogin(server, "resetuser@example.com", "testpassword123", csrfToken, cookie)

	// Request password reset
	resetPayload := map[string]interface{}{"email": "resetuser@example.com"}
	resetBody, _ := json.Marshal(resetPayload)
	resetReq, _ := http.NewRequest("POST", server.URL+"/api/auth/request-reset", bytes.NewReader(resetBody))
	resetReq.Header.Set("Content-Type", "application/json")
	resetReq.Header.Set("X-CSRF-Token", csrfToken)
	resetReq.Header.Set("Cookie", cookie)
	resetResp, err := http.DefaultClient.Do(resetReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, resetResp.StatusCode)
	var resetResult map[string]interface{}
	if err := json.NewDecoder(resetResp.Body).Decode(&resetResult); err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, resetResult["message"], "Password reset email sent")

	// Simulate token (would be sent by email in real app)
	// For test, get token from controllers package (if exported) or skip actual reset
}

func TestHealthAndVersionEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()

	// Health
	healthReq, _ := http.NewRequest("GET", server.URL+"/health", nil)
	healthResp, err := http.DefaultClient.Do(healthReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, healthResp.StatusCode)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(healthResp.Body); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK", buf.String())

	// Version
	versionReq, _ := http.NewRequest("GET", server.URL+"/version", nil)
	versionResp, err := http.DefaultClient.Do(versionReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, versionResp.StatusCode)
	var version map[string]interface{}
	if err := json.NewDecoder(versionResp.Body).Decode(&version); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1.0.0", version["version"])
}

func TestSendMail_Mock(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping send mail mock test in CI")
	}
	// ... existing code ...
}
