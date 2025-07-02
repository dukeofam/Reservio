package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store = session.New(session.Config{
	CookieHTTPOnly: true,
	CookieSecure:   true,
	CookieSameSite: "Strict",
})

func SetSession(c *fiber.Ctx, userID uint) {
	sess, _ := Store.Get(c)
	sess.Set("user_id", strconv.Itoa(int(userID)))
	if err := sess.Save(); err != nil {
		// Optionally log or handle error
	}
}

func ClearSession(c *fiber.Ctx) {
	sess, _ := Store.Get(c)
	if err := sess.Destroy(); err != nil {
		// Optionally log or handle error
	}
}
