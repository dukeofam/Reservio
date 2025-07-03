package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"reservio/config"
	"reservio/middleware"
	"reservio/models"
	"reservio/utils"
	"strings"
	"sync"
	"time"

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

func Register(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate email
	if err := utils.ValidateEmail(body.Email); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		}
		return
	}

	// Validate password
	if err := utils.ValidatePassword(body.Password); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		}
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user := models.User{Email: body.Email, Password: string(hash), Role: "parent"}

	if result := config.DB.Create(&user); result.Error != nil {
		// Check for duplicate email error
		if strings.Contains(result.Error.Error(), "duplicate key value") && strings.Contains(result.Error.Error(), "email") {
			utils.RespondWithValidationError(w, http.StatusConflict, utils.NewValidationError(utils.ErrDuplicateEmail, "Email already registered", map[string]interface{}{
				"email": body.Email,
			}))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not create user")
		}
		return
	}

	utils.SetSession(w, r, user.ID)
	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate email
	if err := utils.ValidateEmail(body.Email); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		}
		return
	}

	// Validate password presence
	if err := utils.ValidatePassword(body.Password); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Password is required")
		}
		return
	}

	la := utils.GetLoginAttempt(body.Email)
	if la.Count >= 5 && time.Now().Unix()-la.LastFailed < 300 {
		utils.RespondWithValidationError(w, http.StatusTooManyRequests, utils.NewValidationError("RATE_LIMIT_EXCEEDED", "Too many failed login attempts. Please try again in 5 minutes.", map[string]interface{}{
			"retry_after": 300 - (time.Now().Unix() - la.LastFailed),
		}))
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		utils.IncrementLoginAttempt(body.Email)
		log.Printf("[Login] Invalid credentials for email: %s", body.Email)
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Invalid credentials", nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		utils.IncrementLoginAttempt(body.Email)
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Invalid credentials", nil))
		return
	}

	utils.ResetLoginAttempt(body.Email)
	utils.SetSession(w, r, user.ID)

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Logged in successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	utils.ClearSession(w, r)
	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Logged out successfully",
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	log.Printf("[GetProfile] userID from context: %d", userID)
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		}))
		return
	}

	log.Printf("[GetProfile] user: email=%s, role=%s", user.Email, user.Role)
	user.Password = "" // Don't send password in response

	utils.RespondWithSuccess(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		}))
		return
	}

	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate email if provided
	if body.Email != "" {
		if err := utils.ValidateEmail(body.Email); err != nil {
			if validationErr, ok := err.(utils.ValidationError); ok {
				utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
			}
			return
		}
		user.Email = body.Email
	}

	// Validate password if provided
	if body.Password != "" {
		if err := utils.ValidatePassword(body.Password); err != nil {
			if validationErr, ok := err.(utils.ValidationError); ok {
				utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
			}
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		user.Password = string(hash)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "duplicate key value") && strings.Contains(err.Error(), "email") {
			utils.RespondWithValidationError(w, http.StatusConflict, utils.NewValidationError(utils.ErrDuplicateEmail, "Email already in use", map[string]interface{}{
				"email": body.Email,
			}))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update profile")
		}
		return
	}

	if body.Password != "" {
		utils.InvalidateAllUserSessions(w, r)
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Profile updated successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate email
	if err := utils.ValidateEmail(body.Email); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		}
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"email": body.Email,
		}))
		return
	}

	token := generateResetToken()
	resetTokens.Lock()
	resetTokens.m[token] = user.ID
	resetTokens.Unlock()

	resetLink := "http://localhost:3000/reset-password?token=" + token
	if err := utils.SendMail(user.Email, "Password Reset", "Reset your password: "+resetLink); err != nil {
		log.Printf("[RequestPasswordReset] Failed to send email: %v", err)
		// Don't fail the request if email fails, just log it
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Password reset email sent",
		"email":   body.Email,
	})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate password
	if err := utils.ValidatePassword(body.Password); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		}
		return
	}

	resetTokens.RLock()
	userID, ok := resetTokens.m[body.Token]
	resetTokens.RUnlock()
	if !ok {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError("INVALID_TOKEN", "Invalid or expired token", map[string]interface{}{
			"token": body.Token,
		}))
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		}))
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user.Password = string(hash)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	resetTokens.Lock()
	delete(resetTokens.m, body.Token)
	resetTokens.Unlock()

	utils.InvalidateAllUserSessions(w, r)

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Password reset successful",
	})
}
