package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email          string `gorm:"uniqueIndex" json:"email"`
	Password       string
	Role           string  // "parent" or "admin"
	Children       []Child `gorm:"foreignKey:ParentID"`
	SessionVersion int     `gorm:"default:1"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Phone          string  `json:"phone"`
	ProfilePicture string  `json:"profile_picture"`
}
