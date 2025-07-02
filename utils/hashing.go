package utils

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var Store *session.Store

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		Store = session.New(session.Config{
			Storage: redis.New(redis.Config{
				URL: redisURL,
			}),
			CookieHTTPOnly: true,
			CookieSecure:   true,
			CookieSameSite: "Strict",
			Expiration:     time.Hour, // 1 hour
		})
	} else {
		Store = session.New(session.Config{
			CookieHTTPOnly: true,
			CookieSecure:   true,
			CookieSameSite: "Strict",
			Expiration:     time.Hour, // 1 hour
		})
	}
}

// Brute-force login attempt tracker (in-memory, can be replaced with Redis)
type LoginAttempt struct {
	Count      int
	LastFailed int64
}

var loginAttempts = struct {
	sync.Mutex
	m map[string]LoginAttempt
}{m: make(map[string]LoginAttempt)}

func IncrementLoginAttempt(email string) int {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	la := loginAttempts.m[email]
	la.Count++
	la.LastFailed = time.Now().Unix()
	loginAttempts.m[email] = la
	return la.Count
}

func ResetLoginAttempt(email string) {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	delete(loginAttempts.m, email)
}

func GetLoginAttempt(email string) LoginAttempt {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	return loginAttempts.m[email]
}

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

// RespondWithError sends a JSON error response and logs the error
func RespondWithError(c *fiber.Ctx, code int, message string) error {
	log.Printf("[ERROR] %s", message)
	return c.Status(code).JSON(fiber.Map{"error": message})
}
