package controllers

import (
	"reservio/config"
	"reservio/models"
	"reservio/utils"

	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"

	"reservio/middleware"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var resetTokens = struct {
	sync.RWMutex
	m map[string]uint
}{m: make(map[string]uint)}

func generateResetToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
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
		return utils.RespondWithError(c, 400, "Invalid input")
	}
	if !utils.IsEmailValid(body.Email) {
		return utils.RespondWithError(c, 400, "Invalid email format")
	}
	if !utils.IsPasswordStrong(body.Password) {
		return utils.RespondWithError(c, 400, "Password must be at least 8 characters")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user := models.User{Email: body.Email, Password: string(hash), Role: "parent"}

	if result := config.DB.Create(&user); result.Error != nil {
		return utils.RespondWithError(c, 500, "Could not create user")
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
		return utils.RespondWithError(c, 400, "Invalid input")
	}
	if !utils.IsEmailValid(body.Email) {
		return utils.RespondWithError(c, 400, "Invalid email format")
	}
	if !utils.IsFieldPresent(body.Password) {
		return utils.RespondWithError(c, 400, "Password is required")
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		log.Printf("[Login] Invalid credentials for email: %s", body.Email)
		return utils.RespondWithError(c, 401, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return utils.RespondWithError(c, 401, "Invalid credentials")
	}

	utils.SetSession(c, user.ID)
	_ = middleware.RegenerateCSRFToken(c)
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
		return utils.RespondWithError(c, 404, "User not found")
	}
	user.Password = "" // Do not expose password
	return c.JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondWithError(c, 404, "User not found")
	}
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, "Invalid input")
	}
	if body.Email != "" && !utils.IsEmailValid(body.Email) {
		return utils.RespondWithError(c, 400, "Invalid email format")
	}
	if body.Password != "" && !utils.IsPasswordStrong(body.Password) {
		return utils.RespondWithError(c, 400, "Password must be at least 8 characters")
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		user.Password = string(hash)
	}
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.RespondWithError(c, 500, "Failed to update profile")
	}
	if body.Password != "" {
		utils.InvalidateAllUserSessions(c)
	}
	return c.JSON(fiber.Map{"message": "Profile updated"})
}

func RequestPasswordReset(c *fiber.Ctx) error {
	type Req struct {
		Email string `json:"email"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, "Invalid input")
	}
	if !utils.IsEmailValid(body.Email) {
		return utils.RespondWithError(c, 400, "Invalid email format")
	}
	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return utils.RespondWithError(c, 404, "User not found")
	}
	token := generateResetToken()
	resetTokens.Lock()
	resetTokens.m[token] = user.ID
	resetTokens.Unlock()
	resetLink := "http://localhost:3000/reset-password?token=" + token
	if err := utils.SendMail(user.Email, "Password Reset", "Reset your password: "+resetLink); err != nil {
		log.Printf("[RequestPasswordReset] Failed to send email: %v", err)
	}
	return c.JSON(fiber.Map{"message": "Password reset email sent"})
}

func ResetPassword(c *fiber.Ctx) error {
	type Req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, "Invalid input")
	}
	if !utils.IsPasswordStrong(body.Password) {
		return utils.RespondWithError(c, 400, "Password must be at least 8 characters")
	}
	resetTokens.RLock()
	userID, ok := resetTokens.m[body.Token]
	resetTokens.RUnlock()
	if !ok {
		return utils.RespondWithError(c, 400, "Invalid or expired token")
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondWithError(c, 404, "User not found")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user.Password = string(hash)
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.RespondWithError(c, 500, "Failed to reset password")
	}
	resetTokens.Lock()
	delete(resetTokens.m, body.Token)
	resetTokens.Unlock()
	utils.InvalidateAllUserSessions(c)
	_ = middleware.RegenerateCSRFToken(c)
	return c.JSON(fiber.Map{"message": "Password reset successful"})
}
