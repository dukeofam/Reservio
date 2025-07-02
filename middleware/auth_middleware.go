package middleware

import (
	"log"
	"strconv"

	"reservio/config"
	"reservio/models"
	"reservio/utils"

	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := utils.Store.Get(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		idStr := sess.Get("user_id")
		if idStr == nil {
			return c.Status(401).JSON(fiber.Map{"error": "Not logged in"})
		}

		id, _ := strconv.Atoi(idStr.(string))
		c.Locals("user_id", uint(id))
		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			log.Printf("[AdminOnly] Forbidden: user not found (user_id=%d)", userID)
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
		}
		if user.Role != "admin" {
			log.Printf("[AdminOnly] Forbidden: user_id=%d, role=%s", userID, user.Role)
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
		}
		return c.Next()
	}
}
