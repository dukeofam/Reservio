package controllers

import (
	"encoding/json"
	"net/http"
	"reservio/config"
	"reservio/models"

	"github.com/gorilla/mux"
)

func CreateSlot(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Date     string `json:"date"`
		Capacity int    `json:"capacity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	slot := models.Slot{Date: body.Date, Capacity: body.Capacity}
	if err := config.DB.Create(&slot).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create slot"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slot); err != nil {
		// Optionally log or handle the error
	}
}

func ApproveReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var reservation models.Reservation
	if err := config.DB.First(&reservation, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Reservation not found"})
		return
	}
	reservation.Status = "approved"
	if err := config.DB.Save(&reservation).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to approve reservation"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reservation); err != nil {
		// Optionally log or handle the error
	}
}

func RejectReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var reservation models.Reservation
	if err := config.DB.First(&reservation, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Reservation not found"})
		return
	}
	reservation.Status = "rejected"
	if err := config.DB.Save(&reservation).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to reject reservation"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reservation); err != nil {
		// Optionally log or handle the error
	}
}

func GetReservationsByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var reservations []models.Reservation
	if status != "" {
		config.DB.Where("status = ?", status).Find(&reservations)
	} else {
		config.DB.Find(&reservations)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reservations); err != nil {
		// Optionally log or handle the error
	}
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	config.DB.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		// Optionally log or handle the error
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"}); err != nil {
		// Optionally log or handle the error
	}
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	user.Role = body.Role
	if err := config.DB.Save(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user role"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		// Optionally log or handle the error
	}
}
