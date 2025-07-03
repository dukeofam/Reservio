package controllers

import (
	"encoding/json"
	"net/http"
	"reservio/config"
	"reservio/middleware"
	"reservio/models"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func MakeReservation(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ChildID uint `json:"child_id"`
		SlotID  uint `json:"slot_id"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	if body.ChildID == 0 || body.SlotID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "child_id and slot_id are required"})
		return
	}
	reservation := models.Reservation{
		ChildID: body.ChildID,
		SlotID:  body.SlotID,
		Status:  "pending",
	}
	if err := config.DB.Create(&reservation).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Reservation failed"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Reservation requested"}); err != nil {
		zap.L().Error("encode error", zap.Error(err))
	}
}

func GetReservations(w http.ResponseWriter, r *http.Request) {
	var reservations []models.Reservation
	config.DB.Find(&reservations)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reservations)
}

func GetMyReservations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	var reservations []models.Reservation
	config.DB.Joins("JOIN children ON children.id = reservations.child_id").Where("children.parent_id = ?", userID).Find(&reservations)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(reservations)
}

func CancelReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := config.DB.Delete(&models.Reservation{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to cancel reservation"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Reservation cancelled"})
}

func ListSlots(w http.ResponseWriter, r *http.Request) {
	var slots []models.Slot
	config.DB.Find(&slots)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(slots)
}
