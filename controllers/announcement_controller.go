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

// ListAnnouncements returns public announcements ordered by newest first (paginated).
func ListAnnouncements(w http.ResponseWriter, r *http.Request) {
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

	var announcements []models.Announcement
	var total int64

	config.DB.Model(&models.Announcement{}).Count(&total)

	offset := (page - 1) * perPage
	if err := config.DB.Preload("Author").Order("created_at DESC").Offset(offset).Limit(perPage).Find(&announcements).Error; err != nil {
		zap.L().Error("Failed to get announcements", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve announcements")
		return
	}

	// Convert to response format (avoid sending author password)
	var data []map[string]interface{}
	for _, a := range announcements {
		item := map[string]interface{}{
			"id":         a.ID,
			"title":      a.Title,
			"content":    a.Content,
			"created_at": a.CreatedAt,
		}
		if a.Author != nil {
			item["author"] = map[string]interface{}{
				"id":    a.Author.ID,
				"email": a.Author.Email,
			}
		}
		data = append(data, item)
	}

	if data == nil {
		data = []map[string]interface{}{}
	}

	utils.RespondWithPaginatedData(w, data, page, perPage, int(total))
}

// CreateAnnouncement (admin only)
func CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	// Must be admin (middleware.AdminOnly already enforced in route)
	userID, _ := r.Context().Value(middleware.UserIDKey).(uint)
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}
	if !utils.IsFieldPresent(body.Title) {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Title is required", nil))
		return
	}
	if !utils.IsFieldPresent(body.Content) {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Content is required", nil))
		return
	}
	ann := models.Announcement{Title: body.Title, Content: body.Content, AuthorID: userID}
	if err := config.DB.Create(&ann).Error; err != nil {
		zap.L().Error("Failed to create announcement", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create announcement")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Announcement created successfully",
		"announcement": map[string]interface{}{
			"id":         ann.ID,
			"title":      ann.Title,
			"content":    ann.Content,
			"created_at": ann.CreatedAt,
		},
	})
}

// UpdateAnnouncement (admin only)
func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	annID, err := utils.ParseUint(idStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid announcement ID", nil))
		return
	}

	var ann models.Announcement
	if err := config.DB.First(&ann, annID).Error; err != nil {
		utils.RespondWithValidationError(w, http.StatusNotFound, utils.NewValidationError(utils.ErrNotFound, "Announcement not found", nil))
		return
	}

	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid JSON input", nil))
		return
	}

	if utils.IsFieldPresent(body.Title) {
		ann.Title = body.Title
	}
	if utils.IsFieldPresent(body.Content) {
		ann.Content = body.Content
	}

	if err := config.DB.Save(&ann).Error; err != nil {
		zap.L().Error("Failed to update announcement", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update announcement")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message": "Announcement updated successfully",
		"announcement": map[string]interface{}{
			"id":         ann.ID,
			"title":      ann.Title,
			"content":    ann.Content,
			"created_at": ann.CreatedAt,
		},
	})
}

// DeleteAnnouncement (admin only)
func DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	annID, err := utils.ParseUint(idStr)
	if err != nil {
		utils.RespondWithValidationError(w, http.StatusBadRequest, utils.NewValidationError(utils.ErrInvalidInput, "Invalid announcement ID", nil))
		return
	}

	if err := config.DB.Delete(&models.Announcement{}, annID).Error; err != nil {
		zap.L().Error("Failed to delete announcement", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete announcement")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"message":         "Announcement deleted successfully",
		"announcement_id": annID,
	})
}
