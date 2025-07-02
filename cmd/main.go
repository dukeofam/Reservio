package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"reservio/config"
	"reservio/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDatabase()

	config.InitSessionStore()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			accept := c.Get("Accept")
			if accept == "application/json" || c.Path() == "/api" || len(c.Path()) > 4 && c.Path()[:5] == "/api/" {
				return c.Status(code).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(code).SendString(err.Error())
		},
	})
	app.Use(logger.New())
	// Add secure headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})
	// Add compression
	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Encoding", "gzip")
		return c.Next()
	})
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{"error": "Too many requests"})
		},
	}))

	routes.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
