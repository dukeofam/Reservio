package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"reservio/config"
	"reservio/models"
	"reservio/utils"
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
		utils.RespondWithError(w, 400, "Invalid input")
		return
	}
	if !utils.IsEmailValid(body.Email) {
		utils.RespondWithError(w, 400, "Invalid email format")
		return
	}
	if !utils.IsPasswordStrong(body.Password) {
		utils.RespondWithError(w, 400, "Password must be at least 8 characters")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user := models.User{Email: body.Email, Password: string(hash), Role: "parent"}

	if result := config.DB.Create(&user); result.Error != nil {
		utils.RespondWithError(w, 500, "Could not create user")
		return
	}

	utils.SetSession(w, r, user.ID)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "User registered", "user": user.Email})
}

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, 400, "Invalid input")
		return
	}
	if !utils.IsEmailValid(body.Email) {
		utils.RespondWithError(w, 400, "Invalid email format")
		return
	}
	if !utils.IsFieldPresent(body.Password) {
		utils.RespondWithError(w, 400, "Password is required")
		return
	}

	la := utils.GetLoginAttempt(body.Email)
	if la.Count >= 5 && time.Now().Unix()-la.LastFailed < 300 {
		utils.RespondWithError(w, 429, "Too many failed login attempts. Please try again in 5 minutes.")
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		utils.IncrementLoginAttempt(body.Email)
		log.Printf("[Login] Invalid credentials for email: %s", body.Email)
		utils.RespondWithError(w, 401, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		utils.IncrementLoginAttempt(body.Email)
		utils.RespondWithError(w, 401, "Invalid credentials")
		return
	}

	utils.ResetLoginAttempt(body.Email)
	utils.SetSession(w, r, user.ID)
	// CSRF token rotation can be handled if needed
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Logged in", "user": user.Email})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	utils.ClearSession(w, r)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Logged out"})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		utils.RespondWithError(w, 401, "Unauthorized")
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(w, 404, "User not found")
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		utils.RespondWithError(w, 401, "Unauthorized")
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(w, 404, "User not found")
		return
	}
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, 400, "Invalid input")
		return
	}
	if body.Email != "" && !utils.IsEmailValid(body.Email) {
		utils.RespondWithError(w, 400, "Invalid email format")
		return
	}
	if body.Password != "" && !utils.IsPasswordStrong(body.Password) {
		utils.RespondWithError(w, 400, "Password must be at least 8 characters")
		return
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		user.Password = string(hash)
	}
	if err := config.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(w, 500, "Failed to update profile")
		return
	}
	if body.Password != "" {
		utils.InvalidateAllUserSessions(w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Profile updated"})
}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, 400, "Invalid input")
		return
	}
	if !utils.IsEmailValid(body.Email) {
		utils.RespondWithError(w, 400, "Invalid email format")
		return
	}
	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		utils.RespondWithError(w, 404, "User not found")
		return
	}
	token := generateResetToken()
	resetTokens.Lock()
	resetTokens.m[token] = user.ID
	resetTokens.Unlock()
	resetLink := "http://localhost:3000/reset-password?token=" + token
	if err := utils.SendMail(user.Email, "Password Reset", "Reset your password: "+resetLink); err != nil {
		log.Printf("[RequestPasswordReset] Failed to send email: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Password reset email sent"})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, 400, "Invalid input")
		return
	}
	if !utils.IsPasswordStrong(body.Password) {
		utils.RespondWithError(w, 400, "Password must be at least 8 characters")
		return
	}
	resetTokens.RLock()
	userID, ok := resetTokens.m[body.Token]
	resetTokens.RUnlock()
	if !ok {
		utils.RespondWithError(w, 400, "Invalid or expired token")
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(w, 404, "User not found")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	user.Password = string(hash)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(w, 500, "Failed to reset password")
		return
	}
	resetTokens.Lock()
	delete(resetTokens.m, body.Token)
	resetTokens.Unlock()
	utils.InvalidateAllUserSessions(w, r)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Password reset successful"})
}
