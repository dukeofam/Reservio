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

func AddChild(w http.ResponseWriter, r *http.Request) {
	type ChildRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body ChildRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	if body.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Name is required"})
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	child := models.Child{Name: body.Name, Age: body.Age, ParentID: userID}
	if err := config.DB.Create(&child).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create child"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(child); err != nil {
		zap.L().Error("encode error", zap.Error(err))
	}
}

func GetChildren(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	var children []models.Child
	config.DB.Where("parent_id = ?", userID).Find(&children)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(children)
}

func EditChild(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	var child models.Child
	if err := config.DB.Where("parent_id = ? AND id = ?", userID, id).First(&child).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Child not found"})
		return
	}
	type Req struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid input"})
		return
	}
	if body.Name != "" {
		child.Name = body.Name
	}
	if body.Age != 0 {
		child.Age = body.Age
	}
	if err := config.DB.Save(&child).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update child"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(child); err != nil {
		zap.L().Error("encode error", zap.Error(err))
	}
}

func DeleteChild(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	result := config.DB.Where("parent_id = ? AND id = ?", userID, id).Delete(&models.Child{})
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete child"})
		return
	}
	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Child not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Child deleted"})
}
