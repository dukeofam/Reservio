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
)

func TestChildEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	_, cookie = registerAndLogin(server, "childparent+1@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+1@example.com").Update("role", "parent")

	// Add child
	addPayload := map[string]interface{}{"name": "TestChild", "age": 5}
	addBody, _ := json.Marshal(addPayload)
	addReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, err := http.DefaultClient.Do(addReq)
	if err != nil {
		t.Fatal(err)
	}
	if addResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", addResp.StatusCode)
	}
	var child map[string]interface{}
	if err := json.NewDecoder(addResp.Body).Decode(&child); err != nil {
		t.Fatal(err)
	}
	if child["name"] != "TestChild" {
		t.Fatalf("expected name 'TestChild', got %v", child["name"])
	}
}

func TestAddGetEditDeleteChild(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	_, cookie = registerAndLogin(server, "childparent+1@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+1@example.com").Update("role", "parent")

	// Add child
	childPayload := map[string]interface{}{"name": "Alice", "age": 5}
	childBody, _ := json.Marshal(childPayload)
	addReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, err := http.DefaultClient.Do(addReq)
	if err != nil {
		t.Fatal(err)
	}
	if addResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", addResp.StatusCode)
	}
	var added map[string]interface{}
	if err := json.NewDecoder(addResp.Body).Decode(&added); err != nil {
		t.Fatal(err)
	}
	if added["Name"] != "Alice" {
		t.Fatalf("expected name 'Alice', got %v", added["Name"])
	}

	// Get children
	getReq, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}
	var children []map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&children); err != nil {
		t.Fatal(err)
	}
	if len(children) == 0 {
		t.Fatalf("expected at least one child")
	}

	// Edit child
	childID := int(added["ID"].(float64))
	editPayload := map[string]interface{}{"name": "AliceUpdated", "age": 6}
	editBody, _ := json.Marshal(editPayload)
	editReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(editBody))
	editReq.Header.Set("Content-Type", "application/json")
	editReq.Header.Set("X-CSRF-Token", csrfToken)
	editReq.Header.Set("Cookie", cookie)
	editResp, err := http.DefaultClient.Do(editReq)
	if err != nil {
		t.Fatal(err)
	}
	if editResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", editResp.StatusCode)
	}
	var edited map[string]interface{}
	if err := json.NewDecoder(editResp.Body).Decode(&edited); err != nil {
		t.Fatal(err)
	}
	if edited["Name"] != "AliceUpdated" {
		t.Fatalf("expected name 'AliceUpdated', got %v", edited["Name"])
	}

	// Delete child
	delReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/children/"+strconv.Itoa(childID), nil)
	delReq.Header.Set("X-CSRF-Token", csrfToken)
	delReq.Header.Set("Cookie", cookie)
	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatal(err)
	}
	if delResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", delResp.StatusCode)
	}
	var delResult map[string]interface{}
	if err := json.NewDecoder(delResp.Body).Decode(&delResult); err != nil {
		t.Fatal(err)
	}
	if delResult["message"] != "Child deleted" {
		t.Fatalf("expected 'Child deleted', got %v", delResult["message"])
	}
}

func TestAddChild_InvalidData(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	_, cookie = registerAndLogin(server, "childparent+2@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+2@example.com").Update("role", "parent")

	// Missing name
	childPayload := map[string]interface{}{"age": 5}
	childBody, _ := json.Marshal(childPayload)
	addReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	addReq.Header.Set("Content-Type", "application/json")
	addReq.Header.Set("X-CSRF-Token", csrfToken)
	addReq.Header.Set("Cookie", cookie)
	addResp, err := http.DefaultClient.Do(addReq)
	if err != nil {
		t.Fatal(err)
	}
	if addResp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", addResp.StatusCode)
	}
}

func TestDeleteChild_NotFound(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)
	_, cookie = registerAndLogin(server, "childparent+3@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+3@example.com").Update("role", "parent")

	delReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/children/99999", nil)
	delReq.Header.Set("X-CSRF-Token", csrfToken)
	delReq.Header.Set("Cookie", cookie)
	delResp, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatal(err)
	}
	if delResp.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", delResp.StatusCode)
	}
}
