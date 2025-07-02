package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store = session.New()

func SetSession(c *fiber.Ctx, userID uint) {
	sess, _ := Store.Get(c)
	sess.Set("user_id", strconv.Itoa(int(userID)))
	sess.Save()
}

func ClearSession(c *fiber.Ctx) {
	sess, _ := Store.Get(c)
	sess.Destroy()
}
