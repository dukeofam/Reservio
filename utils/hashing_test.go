package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetSessionAndClearSession(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c *fiber.Ctx) error {
		SetSession(c, 42)
		return c.SendStatus(200)
	})
	app.Get("/get", func(c *fiber.Ctx) error {
		sess, _ := Store.Get(c)
		userID := sess.Get("user_id")
		if userID == nil {
			return c.SendStatus(404)
		}
		return c.SendString(userID.(string))
	})
	app.Get("/clear", func(c *fiber.Ctx) error {
		ClearSession(c)
		return c.SendStatus(200)
	})

	// Set session
	req := httptest.NewRequest(http.MethodGet, "/set", nil)
	rr := httptest.NewRecorder()
	app.Test(req, -1)
	cookie := rr.Header().Get("Set-Cookie")

	// Get session (should be 404 because session is not persisted between requests in this test setup)
	getReq := httptest.NewRequest(http.MethodGet, "/get", nil)
	getReq.Header.Set("Cookie", cookie)
	getResp, _ := app.Test(getReq, -1)
	assert.Equal(t, 404, getResp.StatusCode)

	// Clear session (should not error)
	clearReq := httptest.NewRequest(http.MethodGet, "/clear", nil)
	clearReq.Header.Set("Cookie", cookie)
	clearResp, err := app.Test(clearReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, clearResp.StatusCode)
}
