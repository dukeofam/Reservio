package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"reservio/config"
	"reservio/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupTestApp() *fiber.App {
	os.Setenv("DATABASE_URL", "postgres://reservio:reservio@localhost:5432/reservio_test?sslmode=disable")
	os.Setenv("TEST_MODE", "1")
	config.ConnectDatabase()
	cleanupTestDB(config.DB)
	app := fiber.New()
	routes.Setup(app)
	return app
}

func cleanupTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users, children, reservations, slots RESTART IDENTITY CASCADE;")
}

func getCSRFTokenAndCookie(app *fiber.App) (string, string) {
	req := httptest.NewRequest("GET", "/api/slots", nil)
	resp, _ := app.Test(req, -1)
	cookie := resp.Header.Get("Set-Cookie")
	csrfToken := resp.Header.Get("X-CSRF-Token")
	return csrfToken, cookie
}

func registerAndLogin(app *fiber.App, email, password, csrfToken, cookie string) (string, string) {
	payload := map[string]string{"email": email, "password": password}
	body, _ := json.Marshal(payload)
	regReq := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	regReq.Header.Set("Content-Type", "application/json")
	regReq.Header.Set("X-CSRF-Token", csrfToken)
	if cookie != "" {
		regReq.Header.Set("Cookie", cookie)
	}
	app.Test(regReq, -1)
	return csrfToken, cookie
}

func createSlot(app *fiber.App, csrfToken, cookie, date string, capacity int) int {
	slotPayload := map[string]interface{}{"date": date, "capacity": capacity}
	slotBody, _ := json.Marshal(slotPayload)
	slotReq := httptest.NewRequest("POST", "/api/admin/slots", bytes.NewReader(slotBody))
	slotReq.Header.Set("Content-Type", "application/json")
	slotReq.Header.Set("X-CSRF-Token", csrfToken)
	slotReq.Header.Set("Cookie", cookie)
	slotResp, _ := app.Test(slotReq, -1)
	var slotResult map[string]interface{}
	json.NewDecoder(slotResp.Body).Decode(&slotResult)
	return int(slotResult["ID"].(float64))
}

func createChild(app *fiber.App, csrfToken, cookie, name, birthdate string) int {
	childPayload := map[string]interface{}{"name": name, "birthdate": birthdate}
	childBody, _ := json.Marshal(childPayload)
	childReq := httptest.NewRequest("POST", "/api/parent/children", bytes.NewReader(childBody))
	childReq.Header.Set("Content-Type", "application/json")
	childReq.Header.Set("X-CSRF-Token", csrfToken)
	childReq.Header.Set("Cookie", cookie)
	childResp, _ := app.Test(childReq, -1)
	var childResult map[string]interface{}
	json.NewDecoder(childResp.Body).Decode(&childResult)
	return int(childResult["ID"].(float64))
}

func createReservation(app *fiber.App, csrfToken, cookie string, slotID, childID int) int {
	resPayload := map[string]interface{}{"slot_id": slotID, "child_id": childID}
	resBody, _ := json.Marshal(resPayload)
	resReq := httptest.NewRequest("POST", "/api/parent/reserve", bytes.NewReader(resBody))
	resReq.Header.Set("Content-Type", "application/json")
	resReq.Header.Set("X-CSRF-Token", csrfToken)
	resReq.Header.Set("Cookie", cookie)
	app.Test(resReq, -1)
	// Fetch reservations and return the latest one for this child and slot
	listReq := httptest.NewRequest("GET", "/api/parent/reservations", nil)
	listReq.Header.Set("Cookie", cookie)
	listResp, _ := app.Test(listReq, -1)
	var reservations []map[string]interface{}
	json.NewDecoder(listResp.Body).Decode(&reservations)
	for _, r := range reservations {
		if int(r["ChildID"].(float64)) == childID && int(r["SlotID"].(float64)) == slotID {
			return int(r["ID"].(float64))
		}
	}
	return 0 // Not found
}
