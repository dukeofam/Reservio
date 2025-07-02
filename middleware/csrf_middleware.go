package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"reservio/utils"

	"github.com/gofiber/fiber/v2"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func CSRFMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, _ := utils.Store.Get(c)
		token := sess.Get("csrf_token")
		if token == nil {
			token = generateCSRFToken()
			sess.Set("csrf_token", token)
			if err := sess.Save(); err != nil {
				log.Printf("[CSRF] sess.Save error: %v", err)
			}
		}

		if os.Getenv("TEST_MODE") == "1" && c.Method() == fiber.MethodGet {
			c.Set("X-CSRF-Token", token.(string))
		}

		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodDelete {
			requestToken := c.Get("X-CSRF-Token")
			if requestToken == "" {
				requestToken = c.FormValue("csrf_token")
			}
			if requestToken != token {
				return c.Status(403).JSON(fiber.Map{"error": "Invalid CSRF token"})
			}
		}
		c.Locals("csrf_token", token)
		return c.Next()
	}
}
