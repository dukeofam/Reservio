package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reservio/config"
	"reservio/middleware"
	"reservio/models"
	"reservio/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func AddChild(w http.ResponseWriter, r *http.Request) {
	type ChildRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body ChildRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate child data
	if err := utils.ValidateChild(body.Name, body.Age); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid child data")
		}
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	child := models.Child{Name: body.Name, Age: body.Age, ParentID: userID}
	if err := config.DB.Create(&child).Error; err != nil {
		zap.L().Error("Failed to create child", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create child")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Child added successfully",
		"child": map[string]interface{}{
			"id":   child.ID,
			"name": child.Name,
			"age":  child.Age,
		},
	})
}

func GetChildren(w http.ResponseWriter, r *http.Request) {
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

	var children []models.Child
	var total int64

	// Get total count
	config.DB.Model(&models.Child{}).Where("parent_id = ?", userID).Count(&total)

	// Get paginated results
	offset := (page - 1) * perPage
	if err := config.DB.Where("parent_id = ?", userID).Offset(offset).Limit(perPage).Find(&children).Error; err != nil {
		zap.L().Error("Failed to get children", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve children")
		return
	}

	// Convert to response format
	var childrenData []map[string]interface{}
	for _, child := range children {
		childrenData = append(childrenData, map[string]interface{}{
			"id":   child.ID,
			"name": child.Name,
			"age":  child.Age,
		})
	}

	// Ensure we always return an empty array instead of nil
	if childrenData == nil {
		childrenData = []map[string]interface{}{}
	}

	utils.RespondWithPaginatedData(w, childrenData, page, perPage, int(total))
}

func EditChild(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	vars := mux.Vars(r)
	childIDStr := vars["id"]
	childID, err := utils.ParseUint(childIDStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid child ID", map[string]interface{}{
			"child_id": childIDStr,
		}))
		return
	}

	// Validate child ownership using business logic
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateChildOwnership(childID, userID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Child not found")
		}
		return
	}

	var child models.Child
	if err := config.DB.Where("parent_id = ? AND id = ?", userID, childID).First(&child).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrChildNotFound, "Child not found", map[string]interface{}{
			"child_id": childID,
		}))
		return
	}

	type Req struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	// Validate updated data
	if body.Name != "" {
		if err := utils.ValidateChild(body.Name, body.Age); err != nil {
			if validationErr, ok := err.(utils.ValidationError); ok {
				utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid child data")
			}
			return
		}
		child.Name = body.Name
	}
	if body.Age != 0 {
		if err := utils.ValidateChild(child.Name, body.Age); err != nil {
			if validationErr, ok := err.(utils.ValidationError); ok {
				utils.RespondWithValidationError(w, http.StatusBadRequest, validationErr)
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid child data")
			}
			return
		}
		child.Age = body.Age
	}

	if err := config.DB.Save(&child).Error; err != nil {
		zap.L().Error("Failed to update child", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update child")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Child updated successfully",
		"child": map[string]interface{}{
			"id":   child.ID,
			"name": child.Name,
			"age":  child.Age,
		},
	})
}

func DeleteChild(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	vars := mux.Vars(r)
	childIDStr := vars["id"]
	childID, err := utils.ParseUint(childIDStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid child ID", map[string]interface{}{
			"child_id": childIDStr,
		}))
		return
	}

	// Validate child ownership using business logic
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateChildOwnership(childID, userID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Child not found")
		}
		return
	}

	result := config.DB.Where("parent_id = ? AND id = ?", userID, childID).Delete(&models.Child{})
	if result.Error != nil {
		zap.L().Error("Failed to delete child", zap.Error(result.Error))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete child")
		return
	}
	if result.RowsAffected == 0 {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrChildNotFound, "Child not found", map[string]interface{}{
			"child_id": childID,
		}))
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message":  "Child deleted successfully",
		"child_id": childID,
	})
}

// GetChild returns a single child belonging to the authenticated parent
func GetChild(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		utils.RespondWithValidationError(w, http.StatusUnauthorized, utils.NewValidationError(utils.ErrUnauthorized, "Not authenticated", nil))
		return
	}

	vars := mux.Vars(r)
	childIDStr := vars["id"]
	childID, err := utils.ParseUint(childIDStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid child ID", map[string]interface{}{
			"child_id": childIDStr,
		}))
		return
	}

	// Validate ownership
	validator := utils.NewBusinessLogicValidator()
	if err := validator.ValidateChildOwnership(childID, userID); err != nil {
		if validationErr, ok := err.(utils.ValidationError); ok {
			utils.RespondWithValidationError(w, http.StatusNotFound, validationErr)
		} else {
			utils.RespondWithError(w, http.StatusNotFound, "Child not found")
		}
		return
	}

	var child models.Child
	if err := config.DB.First(&child, childID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrChildNotFound, "Child not found", map[string]interface{}{
			"child_id": childID,
		}))
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"child": map[string]interface{}{
			"id":   child.ID,
			"name": child.Name,
			"age":  child.Age,
		},
	})
}

// Helper function to parse uint from string
func parseUint(s string) (uint, error) {
	var result uint
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
