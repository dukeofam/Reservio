package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"reservio/utils"
	"time"

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
		expiry := sess.Get("csrf_token_expiry")
		now := time.Now().Unix()
		if token == nil || expiry == nil || now > expiry.(int64) {
			token = generateCSRFToken()
			sess.Set("csrf_token", token)
			sess.Set("csrf_token_expiry", now+7200) // 2 hours
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
				log.Printf("[CSRF] Invalid CSRF token: got=%s expected=%s", requestToken, token)
				return c.Status(403).JSON(fiber.Map{"error": "Invalid CSRF token"})
			}
		}
		c.Locals("csrf_token", token)
		return c.Next()
	}
}

func RegenerateCSRFToken(c *fiber.Ctx) error {
	sess, _ := utils.Store.Get(c)
	token := generateCSRFToken()
	sess.Set("csrf_token", token)
	sess.Set("csrf_token_expiry", time.Now().Unix()+7200) // 2 hours
	if err := sess.Save(); err != nil {
		log.Printf("[CSRF] sess.Save error: %v", err)
		return err
	}
	c.Locals("csrf_token", token)
	return nil
}
