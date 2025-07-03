package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reservio/config"
	"reservio/models"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChildEndpoints(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie := registerAndLogin(server, "childparent+1@example.com", "testpassword123", initToken, initCookie)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)

	// Test add child
	childPayload := map[string]interface{}{"name": "TestChild", "age": 5}
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
	assert.Equal(t, "TestChild", child["name"])
	assert.Equal(t, float64(5), child["age"])

	// Test get children
	getReq, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, getResp.StatusCode)

	var getResult map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&getResult); err != nil {
		t.Fatal(err)
	}

	children := getResult["data"].([]interface{})
	assert.Equal(t, 1, len(children))

	childData := children[0].(map[string]interface{})
	assert.Equal(t, "TestChild", childData["name"])
	assert.Equal(t, float64(5), childData["age"])
}

func TestAddGetEditDeleteChild(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie := registerAndLogin(server, "childparent+1@example.com", "testpassword123", initToken, initCookie)
	fmt.Printf("csrfToken: %s, cookie: %s\n", csrfToken, cookie)

	// Add child
	childPayload := map[string]interface{}{"name": "Alice", "age": 7}
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
	assert.Equal(t, "Alice", child["name"])
	assert.Equal(t, float64(7), child["age"])

	// Get children
	getReq, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, getResp.StatusCode)

	var getResult map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&getResult); err != nil {
		t.Fatal(err)
	}

	children := getResult["data"].([]interface{})
	assert.Equal(t, 1, len(children))

	childData := children[0].(map[string]interface{})
	assert.Equal(t, "Alice", childData["name"])
	assert.Equal(t, float64(7), childData["age"])

	// Edit child
	editPayload := map[string]interface{}{"name": "Alice Updated", "age": 8}
	editBody, _ := json.Marshal(editPayload)
	editReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(editBody))
	editReq.Header.Set("Content-Type", "application/json")
	editReq.Header.Set("X-CSRF-Token", csrfToken)
	editReq.Header.Set("Cookie", cookie)
	editResp, err := http.DefaultClient.Do(editReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, editResp.StatusCode)

	var editResult map[string]interface{}
	if err := json.NewDecoder(editResp.Body).Decode(&editResult); err != nil {
		t.Fatal(err)
	}
	editedChild := editResult["child"].(map[string]interface{})
	assert.Equal(t, "Alice Updated", editedChild["name"])
	assert.Equal(t, float64(8), editedChild["age"])

	// Delete child
	deleteReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/children/"+strconv.Itoa(childID), nil)
	deleteReq.Header.Set("X-CSRF-Token", csrfToken)
	deleteReq.Header.Set("Cookie", cookie)
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, deleteResp.StatusCode)

	var deleteResult map[string]interface{}
	if err := json.NewDecoder(deleteResp.Body).Decode(&deleteResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Child deleted successfully", deleteResult["message"])

	// Verify child is deleted
	getReq2, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq2.Header.Set("Cookie", cookie)
	getResp2, err := http.DefaultClient.Do(getReq2)
	assert.NoError(t, err)
	assert.Equal(t, 200, getResp2.StatusCode)

	var getResult2 map[string]interface{}
	if err := json.NewDecoder(getResp2.Body).Decode(&getResult2); err != nil {
		t.Fatal(err)
	}

	children2 := getResult2["data"].([]interface{})
	assert.Equal(t, 0, len(children2))
}

func TestAddChild_InvalidData(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie = registerAndLogin(server, "childparent+2@example.com", "testpassword123", csrfToken, cookie)

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
	csrfToken, cookie = registerAndLogin(server, "childparent+3@example.com", "testpassword123", csrfToken, cookie)

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

func TestChildValidation(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie := registerAndLogin(server, "childparent+4@example.com", "testpassword123", initToken, initCookie)

	// Test empty name
	emptyNamePayload := map[string]interface{}{"name": "", "age": 5}
	emptyNameBody, _ := json.Marshal(emptyNamePayload)
	emptyNameReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(emptyNameBody))
	emptyNameReq.Header.Set("Content-Type", "application/json")
	emptyNameReq.Header.Set("X-CSRF-Token", csrfToken)
	emptyNameReq.Header.Set("Cookie", cookie)
	emptyNameResp, err := http.DefaultClient.Do(emptyNameReq)
	assert.NoError(t, err)
	assert.Equal(t, 400, emptyNameResp.StatusCode)

	var emptyNameResult map[string]interface{}
	if err := json.NewDecoder(emptyNameResp.Body).Decode(&emptyNameResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Child name is required", emptyNameResult["error"])

	// Test invalid age
	invalidAgePayload := map[string]interface{}{"name": "TestChild", "age": 25}
	invalidAgeBody, _ := json.Marshal(invalidAgePayload)
	invalidAgeReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(invalidAgeBody))
	invalidAgeReq.Header.Set("Content-Type", "application/json")
	invalidAgeReq.Header.Set("X-CSRF-Token", csrfToken)
	invalidAgeReq.Header.Set("Cookie", cookie)
	invalidAgeResp, err := http.DefaultClient.Do(invalidAgeReq)
	assert.NoError(t, err)
	assert.Equal(t, 400, invalidAgeResp.StatusCode)

	var invalidAgeResult map[string]interface{}
	if err := json.NewDecoder(invalidAgeResp.Body).Decode(&invalidAgeResult); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Child age must be between 0 and 18", invalidAgeResult["error"])

	// Test valid child
	validPayload := map[string]interface{}{"name": "ValidChild", "age": 10}
	validBody, _ := json.Marshal(validPayload)
	validReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(validBody))
	validReq.Header.Set("Content-Type", "application/json")
	validReq.Header.Set("X-CSRF-Token", csrfToken)
	validReq.Header.Set("Cookie", cookie)
	validResp, err := http.DefaultClient.Do(validReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, validResp.StatusCode)
}

func TestChildAuthorization(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	initToken, initCookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie := registerAndLogin(server, "childparent+5@example.com", "testpassword123", initToken, initCookie)

	// Create a child for the first parent
	childPayload := map[string]interface{}{"name": "FirstChild", "age": 5}
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

	// Create a second parent and try to access the first parent's child
	secondParentPayload := map[string]string{"email": "secondparent@example.com", "password": "testpassword123"}
	secondParentBody, _ := json.Marshal(secondParentPayload)
	secondParentReq, _ := http.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(secondParentBody))
	secondParentReq.Header.Set("Content-Type", "application/json")
	secondParentReq.Header.Set("X-CSRF-Token", csrfToken)
	secondParentReq.Header.Set("Cookie", cookie)
	_, err = http.DefaultClient.Do(secondParentReq)
	assert.NoError(t, err)

	// Login as second parent
	secondParentLoginPayload := map[string]string{"email": "secondparent@example.com", "password": "testpassword123"}
	secondParentLoginBody, _ := json.Marshal(secondParentLoginPayload)
	secondParentLoginReq, _ := http.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(secondParentLoginBody))
	secondParentLoginReq.Header.Set("Content-Type", "application/json")
	secondParentLoginReq.Header.Set("X-CSRF-Token", csrfToken)
	secondParentLoginReq.Header.Set("Cookie", cookie)
	secondParentLoginResp, err := http.DefaultClient.Do(secondParentLoginReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, secondParentLoginResp.StatusCode)

	// Get cookies from login response
	var secondParentCookies []string
	for _, c := range secondParentLoginResp.Cookies() {
		if c.Name == "session" {
			secondParentCookies = append(secondParentCookies, c.Name+"="+c.Value)
		}
	}
	secondParentCookie := strings.Join(secondParentCookies, "; ")

	// Try to edit the first parent's child
	editPayload := map[string]interface{}{"name": "UnauthorizedEdit", "age": 6}
	editBody, _ := json.Marshal(editPayload)
	editReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(editBody))
	editReq.Header.Set("Content-Type", "application/json")
	editReq.Header.Set("X-CSRF-Token", csrfToken)
	editReq.Header.Set("Cookie", secondParentCookie)
	editResp, err := http.DefaultClient.Do(editReq)
	assert.NoError(t, err)
	assert.Equal(t, 404, editResp.StatusCode) // Should fail due to authorization

	// Try to delete the first parent's child
	deleteReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/children/"+strconv.Itoa(childID), nil)
	deleteReq.Header.Set("X-CSRF-Token", csrfToken)
	deleteReq.Header.Set("Cookie", secondParentCookie)
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	assert.NoError(t, err)
	assert.Equal(t, 404, deleteResp.StatusCode) // Should fail due to authorization
}

func TestChildUnauthorizedAccess(t *testing.T) {
	server := setupTestApp()
	defer server.Close()

	// Test child endpoints without authentication
	// Now that middleware is reordered, Protected runs before CSRF
	testCases := []struct {
		endpoint string
		method   string
	}{
		{"/api/parent/children", "GET"},
		{"/api/parent/children", "POST"}, // Test POST without authentication
	}

	for _, tc := range testCases {
		var req *http.Request
		if tc.method == "POST" {
			// For POST requests, send some data
			payload := map[string]interface{}{"name": "TestChild", "age": 5}
			body, _ := json.Marshal(payload)
			req, _ = http.NewRequest(tc.method, server.URL+tc.endpoint, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(tc.method, server.URL+tc.endpoint, nil)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		// Application returns 401 for unauthorized access to child endpoints
		if resp.StatusCode != 401 {
			t.Fatalf("expected 401 for endpoint %s with method %s, got %d", tc.endpoint, tc.method, resp.StatusCode)
		}
	}
}

func TestChildEditValidation(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie = registerAndLogin(server, "childparent+6@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+6@example.com").Update("role", "parent")

	// Create a child first
	childPayload := map[string]interface{}{"name": "EditTestChild", "age": 4}
	childBody, _ := json.Marshal(childPayload)
	childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	childReq.Header.Set("Cookie", cookie)
	childResp, err := http.DefaultClient.Do(childReq)
	if err != nil {
		t.Fatal(err)
	}
	if childResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", childResp.StatusCode)
	}

	var childResult map[string]interface{}
	if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
		t.Fatal(err)
	}
	child := childResult["child"].(map[string]interface{})
	childID := int(child["id"].(float64))

	// Test edit with invalid data
	invalidEditPayload := map[string]interface{}{"name": "", "age": -5}
	invalidEditBody, _ := json.Marshal(invalidEditPayload)
	invalidEditReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(invalidEditBody))
	invalidEditReq.Header.Set("Content-Type", "application/json")
	invalidEditReq.Header.Set("X-CSRF-Token", csrfToken)
	invalidEditReq.Header.Set("Cookie", cookie)
	invalidEditResp, err := http.DefaultClient.Do(invalidEditReq)
	if err != nil {
		t.Fatal(err)
	}
	// Should fail due to invalid age validation
	if invalidEditResp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", invalidEditResp.StatusCode)
	}

	// Test edit with missing name
	missingNameEditPayload := map[string]interface{}{"age": 4}
	missingNameEditBody, _ := json.Marshal(missingNameEditPayload)
	missingNameEditReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childID), bytes.NewReader(missingNameEditBody))
	missingNameEditReq.Header.Set("Content-Type", "application/json")
	missingNameEditReq.Header.Set("X-CSRF-Token", csrfToken)
	missingNameEditReq.Header.Set("Cookie", cookie)
	missingNameEditResp, err := http.DefaultClient.Do(missingNameEditReq)
	if err != nil {
		t.Fatal(err)
	}
	// Should succeed with valid age only
	if missingNameEditResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", missingNameEditResp.StatusCode)
	}
}

func TestChildMultipleChildren(t *testing.T) {
	server := setupTestApp()
	defer server.Close()
	csrfToken, cookie := getCSRFTokenAndCookie(server)
	csrfToken, cookie = registerAndLogin(server, "childparent+7@example.com", "testpassword123", csrfToken, cookie)

	// Ensure user is parent in DB
	config.DB.Model(&models.User{}).Where("email = ?", "childparent+7@example.com").Update("role", "parent")

	// Create multiple children
	children := []map[string]interface{}{
		{"name": "Child1", "age": 6},
		{"name": "Child2", "age": 5},
		{"name": "Child3", "age": 4},
	}

	var childIDs []int
	for _, child := range children {
		childBody, _ := json.Marshal(child)
		childReq, _ := http.NewRequest("POST", server.URL+"/api/parent/children", bytes.NewReader(childBody))
		childReq.Header.Set("Content-Type", "application/json")
		childReq.Header.Set("X-CSRF-Token", csrfToken)
		childReq.Header.Set("Cookie", cookie)
		childResp, err := http.DefaultClient.Do(childReq)
		if err != nil {
			t.Fatal(err)
		}
		if childResp.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", childResp.StatusCode)
		}

		var childResult map[string]interface{}
		if err := json.NewDecoder(childResp.Body).Decode(&childResult); err != nil {
			t.Fatal(err)
		}
		child := childResult["child"].(map[string]interface{})
		childIDs = append(childIDs, int(child["id"].(float64)))
	}

	// Get all children
	getReq, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", getResp.StatusCode)
	}

	var allChildrenResult map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&allChildrenResult); err != nil {
		t.Fatal(err)
	}
	allChildren := allChildrenResult["data"].([]interface{})
	if len(allChildren) != 3 {
		t.Fatalf("expected 3 children, got %d", len(allChildren))
	}

	// Verify all children are present
	foundNames := make(map[string]bool)
	for _, childInterface := range allChildren {
		child := childInterface.(map[string]interface{})
		foundNames[child["name"].(string)] = true
	}
	for _, child := range children {
		if !foundNames[child["name"].(string)] {
			t.Fatalf("child %s not found in list", child["name"])
		}
	}

	// Edit one child
	editPayload := map[string]interface{}{"name": "Child1Updated", "age": 6}
	editBody, _ := json.Marshal(editPayload)
	editReq, _ := http.NewRequest("PUT", server.URL+"/api/parent/children/"+strconv.Itoa(childIDs[0]), bytes.NewReader(editBody))
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

	var editResult map[string]interface{}
	if err := json.NewDecoder(editResp.Body).Decode(&editResult); err != nil {
		t.Fatal(err)
	}
	child := editResult["child"].(map[string]interface{})
	if child["name"] != "Child1Updated" {
		t.Fatalf("expected name 'Child1Updated', got %v", child["name"])
	}

	// Delete one child
	deleteReq, _ := http.NewRequest("DELETE", server.URL+"/api/parent/children/"+strconv.Itoa(childIDs[1]), nil)
	deleteReq.Header.Set("X-CSRF-Token", csrfToken)
	deleteReq.Header.Set("Cookie", cookie)
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatal(err)
	}
	if deleteResp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", deleteResp.StatusCode)
	}

	// Verify only 2 children remain
	getReq2, _ := http.NewRequest("GET", server.URL+"/api/parent/children", nil)
	getReq2.Header.Set("Cookie", cookie)
	getResp2, err := http.DefaultClient.Do(getReq2)
	if err != nil {
		t.Fatal(err)
	}
	if getResp2.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", getResp2.StatusCode)
	}

	var remainingChildrenResult map[string]interface{}
	if err := json.NewDecoder(getResp2.Body).Decode(&remainingChildrenResult); err != nil {
		t.Fatal(err)
	}
	remainingChildren := remainingChildrenResult["data"].([]interface{})
	if len(remainingChildren) != 2 {
		t.Fatalf("expected 2 children after deletion, got %d", len(remainingChildren))
	}
}
