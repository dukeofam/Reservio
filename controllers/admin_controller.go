package controllers

import (
	"strconv"

	"reservio/config"
	"reservio/models"

	"github.com/gofiber/fiber/v2"
)

func CreateSlot(c *fiber.Ctx) error {
	type SlotRequest struct {
		Date     string `json:"date"`
		Capacity int    `json:"capacity"`
	}

	var body SlotRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	slot := models.Slot{Date: body.Date, Capacity: body.Capacity}
	if err := config.DB.Create(&slot).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create slot"})
	}

	return c.JSON(slot)
}

func ApproveReservation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	config.DB.Model(&models.Reservation{}).Where("id = ?", id).Update("status", "approved")
	return c.JSON(fiber.Map{"message": "Reservation approved"})
}

func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	for i := range users {
		users[i].Password = "" // Hide password
	}
	return c.JSON(users)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}

func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	type Req struct {
		Role string `json:"role"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	user.Role = body.Role
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update role"})
	}
	return c.JSON(fiber.Map{"message": "User role updated"})
}

func RejectReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	config.DB.Model(&models.Reservation{}).Where("id = ?", id).Update("status", "rejected")
	return c.JSON(fiber.Map{"message": "Reservation rejected"})
}

func GetReservationsByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	var reservations []models.Reservation
	if status != "" {
		config.DB.Where("status = ?", status).Find(&reservations)
	} else {
		config.DB.Find(&reservations)
	}
	return c.JSON(reservations)
}
