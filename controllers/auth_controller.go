package controllers

import (
	"reservio/config"
	"reservio/models"
	"reservio/utils"

	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var resetTokens = struct {
	sync.RWMutex
	m map[string]uint
}{m: make(map[string]uint)}

func generateResetToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if !utils.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if !utils.IsPasswordStrong(body.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user := models.User{Email: body.Email, Password: string(hash), Role: "parent"}

	if result := config.DB.Create(&user); result.Error != nil {
		log.Printf("[Register] DB error: %v", result.Error)
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	utils.SetSession(c, user.ID)
	return c.JSON(fiber.Map{"message": "User registered", "user": user.Email})
}

func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if !utils.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if !utils.IsFieldPresent(body.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		log.Printf("[Login] Invalid credentials for email: %s", body.Email)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	utils.SetSession(c, user.ID)
	return c.JSON(fiber.Map{"message": "Logged in", "user": user.Email})
}

func Logout(c *fiber.Ctx) error {
	utils.ClearSession(c)
	return c.JSON(fiber.Map{"message": "Logged out"})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Printf("[GetProfile] User not found: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	user.Password = "" // Do not expose password
	return c.JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if body.Email != "" && !utils.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	if body.Password != "" && !utils.IsPasswordStrong(body.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		user.Password = string(hash)
	}
	if err := config.DB.Save(&user).Error; err != nil {
		log.Printf("[UpdateProfile] DB error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update profile"})
	}
	return c.JSON(fiber.Map{"message": "Profile updated"})
}

func RequestPasswordReset(c *fiber.Ctx) error {
	type Req struct {
		Email string `json:"email"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if !utils.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}
	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		log.Printf("[RequestPasswordReset] User not found: %s", body.Email)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	token := generateResetToken()
	resetTokens.Lock()
	resetTokens.m[token] = user.ID
	resetTokens.Unlock()
	resetLink := "http://localhost:3000/reset-password?token=" + token
	utils.SendMail(user.Email, "Password Reset", "Reset your password: "+resetLink)
	return c.JSON(fiber.Map{"message": "Password reset email sent"})
}

func ResetPassword(c *fiber.Ctx) error {
	type Req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if !utils.IsPasswordStrong(body.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}
	resetTokens.RLock()
	userID, ok := resetTokens.m[body.Token]
	resetTokens.RUnlock()
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Printf("[ResetPassword] User not found for token: %s", body.Token)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user.Password = string(hash)
	config.DB.Save(&user)
	resetTokens.Lock()
	delete(resetTokens.m, body.Token)
	resetTokens.Unlock()
	return c.JSON(fiber.Map{"message": "Password reset successful"})
}
