package controllers

import (
	"reservio/config"
	"reservio/models"

	"github.com/gofiber/fiber/v2"
)

func AddChild(c *fiber.Ctx) error {
	type ChildRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body ChildRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}
	userID := c.Locals("user_id").(uint)
	child := models.Child{Name: body.Name, Age: body.Age, ParentID: userID}
	if err := config.DB.Create(&child).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create child"})
	}
	return c.JSON(child)
}

func GetChildren(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var children []models.Child
	config.DB.Where("parent_id = ?", userID).Find(&children)
	return c.JSON(children)
}

func EditChild(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")
	var child models.Child
	if err := config.DB.Where("parent_id = ? AND id = ?", userID, id).First(&child).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Child not found"})
	}
	type Req struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if body.Name != "" {
		child.Name = body.Name
	}
	if body.Age != 0 {
		child.Age = body.Age
	}
	if err := config.DB.Save(&child).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update child"})
	}
	return c.JSON(child)
}

func DeleteChild(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")
	result := config.DB.Where("parent_id = ? AND id = ?", userID, id).Delete(&models.Child{})
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete child"})
	}
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Child not found"})
	}
	return c.JSON(fiber.Map{"message": "Child deleted"})
}
