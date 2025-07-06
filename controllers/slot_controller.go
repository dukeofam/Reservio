package controllers

import (
	"net/http"
	"reservio/config"
	"reservio/models"
	"reservio/utils"
)

// ListSlotsCalendar returns slots grouped by date with availability counts.
func ListSlotsCalendar(w http.ResponseWriter, r *http.Request) {
	var slots []models.Slot
	if err := config.DB.Find(&slots).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve slots")
		return
	}

	// Build map date->slots
	availabilityMap := make(map[string][]map[string]interface{})
	for _, slot := range slots {
		var approved int64
		config.DB.Table("reservations").Where("slot_id = ? AND status IN ('approved','pending')", slot.ID).Count(&approved)
		remaining := slot.Capacity - int(approved)
		item := map[string]interface{}{
			"id":        slot.ID,
			"date":      slot.Date,
			"capacity":  slot.Capacity,
			"remaining": remaining,
		}
		availabilityMap[slot.Date] = append(availabilityMap[slot.Date], item)
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"calendar": availabilityMap,
	})
}
