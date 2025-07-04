package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reservio/config"
	"reservio/routes"
	"strings"

	"gorm.io/gorm"
)

func setupTestApp() *httptest.Server {
	_ = os.Setenv("DATABASE_URL", "postgres://reservio:reservio@localhost:5432/reservio_test?sslmode=disable")
	_ = os.Setenv("TEST_MODE", "1")
	_ = os.Setenv("SESSION_SECRET", "test-secret-key")
	config.ConnectDatabase()
	config.InitSessionStore()
	cleanupTestDB(config.DB)
	router := routes.SetupRouter()
	return httptest.NewServer(router)
}

func cleanupTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users, children, reservations, slots RESTART IDENTITY CASCADE;")
}

func getCSRFTokenAndCookie(server *httptest.Server) (string, string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+"/api/slots", nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	raw := resp.Header.Get("Set-Cookie")
	var cookie string
	if raw != "" {
		parts := strings.SplitN(raw, ";", 2)
		cookie = strings.TrimSpace(parts[0])
	}
	// Note: This function is only used for initial CSRF token before authentication
	// The actual CSRF token for authenticated requests should come from registerAndLogin
	csrfToken := resp.Header.Get("X-CSRF-Token")
	return csrfToken, cookie
}

func registerAndLogin(server *httptest.Server, email, password, csrfToken, cookie string) (string, string) {
	client := &http.Client{}
	payload := map[string]string{"email": email, "password": password}
	body, _ := json.Marshal(payload)

	// Register
	regReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		regReq.Header.Set("Cookie", cookie)
	}
	regResp, err := client.Do(regReq)
	if err != nil {
		panic(err)
	}
	for _, c := range regResp.Cookies() {
		if c.Name == "session" {
			cookie = c.Name + "=" + c.Value
		}
	}
	_ = regResp.Body.Close()

	// Login
	loginReq, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("X-CSRF-Token", csrfToken)
	loginReq.Header.Set("Cookie", cookie)
	loginResp, err := client.Do(loginReq)
	if err != nil {
		panic(err)
	}
	if loginResp.StatusCode != 200 {
		var bodyBytes bytes.Buffer
		_, _ = bodyBytes.ReadFrom(loginResp.Body)
		panic(fmt.Sprintf("login failed: status %d, body %s", loginResp.StatusCode, bodyBytes.String()))
	}
	// Update cookie from login response
	for _, c := range loginResp.Cookies() {
		if c.Name == "session" {
			cookie = c.Name + "=" + c.Value
		}
	}
	_ = loginResp.Body.Close()

	// ⬇️ NEW: Get fresh CSRF after login from a protected endpoint
	getReq, _ := http.NewRequest("GET", server.URL+"/api/user/profile", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := client.Do(getReq)
	if err != nil {
		panic(err)
	}
	defer func() { _ = getResp.Body.Close() }()

	// Get fresh CSRF token from the response
	freshCSRFToken := getResp.Header.Get("X-CSRF-Token")

	// Update cookie with any new Set-Cookie header (session now contains csrf_token)
	for _, c := range getResp.Cookies() {
		if c.Name == "session" {
			cookie = c.Name + "=" + c.Value
		}
	}
	return freshCSRFToken, cookie
}

func createSlot(server *httptest.Server, csrfToken, cookie, date string, capacity int) int {
	client := &http.Client{}

	slotPayload := map[string]interface{}{"date": date, "capacity": capacity}
	slotBody, _ := json.Marshal(slotPayload)
	slotReq, _ := http.NewRequest("POST", server.URL+"/api/admin/slots", bytes.NewReader(slotBody))
	slotReq.Header.Set("Content-Type", "application/json")
	slotReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		slotReq.Header.Set("Cookie", cookie)
	}
	slotResp, err := client.Do(slotReq)
	if err != nil {
		panic(err)
	}
	defer func() { _ = slotResp.Body.Close() }()
	if slotResp.StatusCode != 200 {
		var bodyBytes bytes.Buffer
		_, _ = bodyBytes.ReadFrom(slotResp.Body)
		panic(fmt.Sprintf("createSlot failed: status %d, body: %s", slotResp.StatusCode, bodyBytes.String()))
	}
	var slotResult map[string]interface{}
	if err := json.NewDecoder(slotResp.Body).Decode(&slotResult); err != nil {
		panic(err)
	}
	slot := slotResult["slot"].(map[string]interface{})
	idVal, ok := slot["id"]
	if !ok || idVal == nil {
		panic("createSlot: expected slot ID in response, got: " + fmt.Sprintf("%v", slotResult))
	}
	return int(idVal.(float64))
}

func createChild(server *httptest.Server, csrfToken, cookie, name string, age int) int {
	client := &http.Client{}
	childPayload := map[string]interface{}{"name": name, "age": age}
	childBody, _ := json.Marshal(childPayload)
	childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		childReq.Header.Set("Cookie", cookie)
	}
	childResp, err := client.Do(childReq)
	if err != nil {
		panic(err)
	}
	defer func() { _ = childResp.Body.Close() }()
	if childResp.StatusCode != 200 {
		var bodyBytes bytes.Buffer
		_, _ = bodyBytes.ReadFrom(childResp.Body)
		panic(fmt.Sprintf("createChild: expected 200, got %d, body: %s", childResp.StatusCode, bodyBytes.String()))
	}
	var childResult map[string]interface{}
	if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
		panic(err)
	}
	child := childResult["child"].(map[string]interface{})
	idVal, ok := child["id"]
	if !ok || idVal == nil {
		panic("createChild: expected child ID in response, got: " + fmt.Sprintf("%v", childResult))
	}
	return int(idVal.(float64))
}

func createReservation(server *httptest.Server, csrfToken, cookie string, slotID, childID int) int {
	client := &http.Client{}
	resPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	resBody, _ := json.Marshal(resPayload)
	resReq, _ := http.NewRequest("POST", server.URL+"/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		resReq.Header.Set("Cookie", cookie)
	}
	if _, err := client.Do(resReq); err != nil {
		panic(err)
	}
	// Fetch reservations and return the latest one for this child and slot
	listReq, _ := http.NewRequest("GET", server.URL+"/api/parent/reservations", nil)
	if cookie != "" {
		listReq.Header.Set("Cookie", cookie)
	}
	listResp, err := client.Do(listReq)
	if err != nil {
		panic(err)
	}
	defer func() { _ = listResp.Body.Close() }()
	var listResult map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listResult); err != nil {
		panic(err)
	}

	reservations := listResult["data"].([]interface{})
	for _, rInterface := range reservations {
		r := rInterface.(map[string]interface{})
		if int(r["child_id"].(float64)) == childID && int(r["slot_id"].(float64)) == slotID {
			return int(r["id"].(float64))
		}
	}
	return 0 // Not found
}
