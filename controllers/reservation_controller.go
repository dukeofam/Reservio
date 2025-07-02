package controllers

import (
	"reservio/config"
	"reservio/models"

	"github.com/gofiber/fiber/v2"
)

func MakeReservation(c *fiber.Ctx) error {
	type Req struct {
		ChildID uint `json:"child_id"`
		SlotID  uint `json:"slot_id"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if body.ChildID == 0 || body.SlotID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "child_id and slot_id are required"})
	}

	reservation := models.Reservation{
		ChildID: body.ChildID,
		SlotID:  body.SlotID,
		Status:  "pending",
	}

	if err := config.DB.Create(&reservation).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Reservation failed"})
	}

	return c.JSON(fiber.Map{"message": "Reservation requested"})
}

func GetReservations(c *fiber.Ctx) error {
	var reservations []models.Reservation
	config.DB.Find(&reservations)
	return c.JSON(reservations)
}

func GetMyReservations(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var reservations []models.Reservation
	config.DB.Joins("JOIN children ON children.id = reservations.child_id").Where("children.parent_id = ?", userID).Find(&reservations)
	return c.JSON(reservations)
}

func CancelReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.Reservation{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel reservation"})
	}
	return c.JSON(fiber.Map{"message": "Reservation cancelled"})
}

func ListSlots(c *fiber.Ctx) error {
	var slots []models.Slot
	config.DB.Find(&slots)
	return c.JSON(slots)
}
