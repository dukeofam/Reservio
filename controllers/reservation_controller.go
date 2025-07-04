package controllers

import (
	"encoding/json"
	"net/http"
	"reservio/config"
	"reservio/middleware"
	"reservio/models"
	"reservio/utils"

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
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate reservation data
	if err := utils.ValidateReservation(body.ChildID, body.SlotID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid reservation data")
		}
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	// Business logic validation
	validator := utils.NewBusinessLogicValidator()

	// Validate child ownership
	if err := validator.ValidateChildOwnership(body.ChildID, userID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Child not found")
		}
		return
	}

	// Validate slot exists
	if err := validator.ValidateSlotExists(body.SlotID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Slot not found")
		}
		return
	}

	// Validate slot capacity
	if err := validator.ValidateSlotCapacity(body.SlotID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusConflict, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusConflict, "Slot is full")
		}
		return
	}

	// Validate no double booking
	if err := validator.ValidateDoubleBooking(body.ChildID, body.SlotID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusConflict, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusConflict, "Child is already booked for this slot")
		}
		return
	}

	reservation := models.Reservation{
		ChildID: body.ChildID,
		SlotID:  body.SlotID,
		Status:  "pending",
	}

	if err := config.DB.Create(&reservation).Error; err != nil {
		zap.L().Error("Failed to create reservation", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create reservation")
		return
	}

	// Get slot availability for response
	availability, _ := validator.GetSlotAvailability(body.SlotID)

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Reservation requested successfully",
		"reservation": map[string]interface{}{
			"id":       reservation.ID,
			"child_id": reservation.ChildID,
			"slot_id":  reservation.SlotID,
			"status":   reservation.Status,
		},
		"slot_availability": availability,
	})
}

func GetReservations(w http.ResponseWriter, r *http.Request) {
	var reservations []models.Reservation
	config.DB.Find(&reservations)
	utils.RespondWithSuccess(w, map[string]interface{}{
		"reservations": reservations,
	})
}

func GetMyReservations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	// Parse pagination parameters
	page, perPage, err := utils.ParsePagination(r.URL.Query().Get("page"), r.URL.Query().Get("per_page"))
	if err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		}
		return
	}

	// Parse status filter
	status := r.URL.Query().Get("status")
	if status != "" {
		if err := utils.ValidateStatus(status); err != nil {
			if validationErr, ok := err.(utils.ValidationError); ok {
				utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid status")
			}
			return
		}
	}

	var reservations []models.Reservation
	var total int64

	// Build query
	query := config.DB.Joins("JOIN children ON children.id = reservations.child_id").Where("children.parent_id = ?", userID)
	if status != "" {
		query = query.Where("reservations.status = ?", status)
	}

	// Get total count
	query.Model(&models.Reservation{}).Count(&total)

	// Get paginated results
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Find(&reservations).Error; err != nil {
		zap.L().Error("Failed to get reservations", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve reservations")
		return
	}

	// Convert to response format
	var reservationsData []map[string]interface{}
	for _, reservation := range reservations {
		reservationData := map[string]interface{}{
			"id":       reservation.ID,
			"child_id": reservation.ChildID,
			"slot_id":  reservation.SlotID,
			"status":   reservation.Status,
		}

		reservationsData = append(reservationsData, reservationData)
	}

	// Ensure we always return an empty array instead of nil
	if reservationsData == nil {
		reservationsData = []map[string]interface{}{}
	}

	utils.RespondWithPaginatedData(w, reservationsData, page, perPage, int(total))
}

func CancelReservation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	vars := mux.Vars(r)
	reservationIDStr := vars["id"]
	reservationID, err := utils.ParseUint(reservationIDStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid reservation ID", map[string]interface{}{
			"reservation_id": reservationIDStr,
		}))
		return
	}

	// Validate reservation ownership using business logic
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateReservationOwnership(reservationID, userID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Reservation not found")
		}
		return
	}

	if err := config.DB.Delete(&models.Reservation{}, reservationID).Error; err != nil {
		zap.L().Error("Failed to cancel reservation", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel reservation")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message":        "Reservation cancelled successfully",
		"reservation_id": reservationID,
	})
}

func ListSlots(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, perPage, err := utils.ParsePagination(r.URL.Query().Get("page"), r.URL.Query().Get("per_page"))
	if err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		}
		return
	}

	var slots []models.Slot
	var total int64

	// Get total count
	config.DB.Model(&models.Slot{}).Count(&total)

	// Get paginated results
	offset := (page - 1) * perPage
	if err := config.DB.Offset(offset).Limit(perPage).Find(&slots).Error; err != nil {
		zap.L().Error("Failed to get slots", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve slots")
		return
	}

	// Add availability information to each slot
	validator := utils.NewBusinessLogicValidator()
	var slotsData []map[string]interface{}
	for _, slot := range slots {
		availability, _ := validator.GetSlotAvailability(slot.ID)

		slotData := map[string]interface{}{
			"id":       slot.ID,
			"date":     slot.Date,
			"capacity": slot.Capacity,
		}

		// Merge availability information
		for k, v := range availability {
			slotData[k] = v
		}

		slotsData = append(slotsData, slotData)
	}

	// Ensure we always return an empty array instead of nil
	if slotsData == nil {
		slotsData = []map[string]interface{}{}
	}

	utils.RespondWithPaginatedData(w, slotsData, page, perPage, int(total))
}

// GetSlot returns detailed information (including availability) about a single slot
func GetSlot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slotIDStr := vars["id"]
	slotID, err := utils.ParseUint(slotIDStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid slot ID", map[string]interface{}{
			"slot_id": slotIDStr,
		}))
		return
	}

	validator := utils.NewBusinessLogicValidator()
	// Ensure slot exists
	if err := validator.ValidateSlotExists(slotID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Slot not found")
		}
		return
	}

	var slot models.Slot
	if err := config.DB.First(&slot, slotID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrSlotNotFound, "Slot not found", map[string]interface{}{
			"slot_id": slotID,
		}))
		return
	}

	availability, _ := validator.GetSlotAvailability(slotID)

	response := map[string]interface{}{
		"slot": map[string]interface{}{
			"id":       slot.ID,
			"date":     slot.Date,
			"capacity": slot.Capacity,
		},
		"availability": availability,
	}

	utils.RespondWithSuccess(w, response)
}
