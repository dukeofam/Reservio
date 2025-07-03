package controllers

import (
	"encoding/json"
	"net/http"
	"reservio/config"
	"reservio/models"
	"reservio/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func CreateSlot(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Date     string `json:"date"`
		Capacity int    `json:"capacity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate slot data
	if err := utils.ValidateSlot(body.Date, body.Capacity); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid slot data")
		}
		return
	}

	slot := models.Slot{Date: body.Date, Capacity: body.Capacity}
	if err := config.DB.Create(&slot).Error; err != nil {
		zap.L().Error("Failed to create slot", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create slot")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Slot created successfully",
		"slot": map[string]interface{}{
			"id":       slot.ID,
			"date":     slot.Date,
			"capacity": slot.Capacity,
		},
	})
}

func ApproveReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reservationID, err := utils.ParseUint(id)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid reservation ID", map[string]interface{}{
			"reservation_id": id,
		}))
		return
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, reservationID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrReservationNotFound, "Reservation not found", map[string]interface{}{
			"reservation_id": reservationID,
		}))
		return
	}

	// Validate status transition
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateReservationStatusTransition(reservation.Status, "approved"); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid status transition")
		}
		return
	}

	reservation.Status = "approved"
	if err := config.DB.Save(&reservation).Error; err != nil {
		zap.L().Error("Failed to approve reservation", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to approve reservation")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Reservation approved successfully",
		"reservation": map[string]interface{}{
			"id":       reservation.ID,
			"child_id": reservation.ChildID,
			"slot_id":  reservation.SlotID,
			"status":   reservation.Status,
		},
	})
}

func RejectReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reservationID, err := utils.ParseUint(id)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid reservation ID", map[string]interface{}{
			"reservation_id": id,
		}))
		return
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, reservationID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrReservationNotFound, "Reservation not found", map[string]interface{}{
			"reservation_id": reservationID,
		}))
		return
	}

	// Validate status transition
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateReservationStatusTransition(reservation.Status, "rejected"); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid status transition")
		}
		return
	}

	reservation.Status = "rejected"
	if err := config.DB.Save(&reservation).Error; err != nil {
		zap.L().Error("Failed to reject reservation", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to reject reservation")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Reservation rejected successfully",
		"reservation": map[string]interface{}{
			"id":       reservation.ID,
			"child_id": reservation.ChildID,
			"slot_id":  reservation.SlotID,
			"status":   reservation.Status,
		},
	})
}

func GetReservationsByStatus(w http.ResponseWriter, r *http.Request) {
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
	query := config.DB.Model(&models.Reservation{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	query.Count(&total)

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

func ListUsers(w http.ResponseWriter, r *http.Request) {
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

	var users []models.User
	var total int64

	// Get total count
	config.DB.Model(&models.User{}).Count(&total)

	// Get paginated results
	offset := (page - 1) * perPage
	if err := config.DB.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		zap.L().Error("Failed to get users", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Convert to response format (exclude passwords)
	var usersData []map[string]interface{}
	for _, user := range users {
		userData := map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		}

		usersData = append(usersData, userData)
	}

	// Ensure we always return an empty array instead of nil
	if usersData == nil {
		usersData = []map[string]interface{}{}
	}

	utils.RespondWithPaginatedData(w, usersData, page, perPage, int(total))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := utils.ParseUint(id)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}))
		return
	}

	result := config.DB.Delete(&models.User{}, userID)
	if result.Error != nil {
		zap.L().Error("Failed to delete user", zap.Error(result.Error))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	if result.RowsAffected == 0 {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		}))
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := utils.ParseUint(id)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}))
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate role
	if err := utils.ValidateRole(body.Role); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid role")
		}
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		}))
		return
	}

	user.Role = body.Role
	if err := config.DB.Save(&user).Error; err != nil {
		zap.L().Error("Failed to update user role", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "User role updated successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
