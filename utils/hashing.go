package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store = session.New(session.Config{
	CookieHTTPOnly: true,
	CookieSecure:   true,
	CookieSameSite: "Strict",
	Expiration:     time.Hour, // 1 hour
})

func SetSession(c *fiber.Ctx, userID uint) {
	sess, _ := Store.Get(c)
	sess.Set("user_id", strconv.Itoa(int(userID)))
	if err := sess.Save(); err != nil {
		log.Printf("[SetSession] sess.Save error: %v", err)
	}
}

func ClearSession(c *fiber.Ctx) {
	sess, _ := Store.Get(c)
	if err := sess.Destroy(); err != nil {
		log.Printf("[ClearSession] sess.Destroy error: %v", err)
	}
}

// InvalidateAllUserSessions destroys the current session and rotates the session ID.
// NOTE: To fully invalidate all sessions for a user across devices, use a distributed session store (e.g., Redis)
// and track sessions by user ID. This implementation only affects the current session.
func InvalidateAllUserSessions(c *fiber.Ctx) {
	sess, _ := Store.Get(c)
	_ = sess.Destroy()
	_ = sess.Regenerate()
}
