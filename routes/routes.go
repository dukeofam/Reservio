package routes

import (
	"reservio/controllers"
	"reservio/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", middleware.CSRFMiddleware())

	auth := api.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Post("/request-reset", controllers.RequestPasswordReset)
	auth.Post("/reset-password", controllers.ResetPassword)

	parent := api.Group("/parent", middleware.Protected())
	parent.Post("/children", controllers.AddChild)
	parent.Get("/children", controllers.GetChildren)
	parent.Post("/reserve", controllers.MakeReservation)
	parent.Get("/reservations", controllers.GetMyReservations)
	parent.Delete("/reservations/:id", controllers.CancelReservation)
	parent.Put("/children/:id", controllers.EditChild)
	parent.Delete("/children/:id", controllers.DeleteChild)

	admin := api.Group("/admin", middleware.Protected(), middleware.AdminOnly())
	admin.Post("/slots", controllers.CreateSlot)
	admin.Put("/approve/:id", controllers.ApproveReservation)
	admin.Put("/reject/:id", controllers.RejectReservation)
	admin.Get("/reservations", controllers.GetReservationsByStatus)
	admin.Get("/users", controllers.ListUsers)
	admin.Delete("/users/:id", controllers.DeleteUser)
	admin.Put("/users/:id/role", controllers.UpdateUserRole)

	user := api.Group("/user", middleware.Protected())
	user.Get("/profile", controllers.GetProfile)
	user.Put("/profile", controllers.UpdateProfile)

	api.Get("/slots", controllers.ListSlots)

	// Add health and version endpoints at root
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"version": "1.0.0", "commit": "dev"})
	})
}
