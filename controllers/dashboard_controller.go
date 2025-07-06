package controllers

import (
	"net/http"
	"reservio/config"
	"reservio/middleware"
	"reservio/models"
	"reservio/utils"
)

// GetDashboardStats returns counts for children, reservations and open slots.
// For admins it returns totals across the system.
// For parents it returns counts scoped to their data.
func GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	var user models.User
	if err := config.DB.Select("role").First(&user, userID).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user role")
		return
	}

	var totalChildren int64
	var totalReservations int64
	var openSlots int64

	if user.Role == "admin" {
		config.DB.Model(&models.Child{}).Count(&totalChildren)
		config.DB.Model(&models.Reservation{}).Count(&totalReservations)
		// open slots: count of slots where reservations < capacity
		config.DB.Model(&models.Slot{}).
			Joins("LEFT JOIN reservations ON reservations.slot_id = slots.id AND reservations.status IN ('approved','pending')").
			Group("slots.id").
			Having("COUNT(reservations.id) < slots.capacity").
			Count(&openSlots)
	} else {
		// Parent
		config.DB.Model(&models.Child{}).Where("parent_id = ?", userID).Count(&totalChildren)
		// Reservations via join on children
		config.DB.Table("reservations").
			Joins("JOIN children ON children.id = reservations.child_id").
			Where("children.parent_id = ?", userID).
			Count(&totalReservations)
		// Open slots = slots with availability (not user specific)
		config.DB.Model(&models.Slot{}).
			Joins("LEFT JOIN reservations ON reservations.slot_id = slots.id AND reservations.status IN ('approved','pending')").
			Group("slots.id").
			Having("COUNT(reservations.id) < slots.capacity").
			Count(&openSlots)
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"total_children":     totalChildren,
		"total_reservations": totalReservations,
		"open_slots":         openSlots,
	})
}
