package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string
	Role     string  // "parent" or "admin"
	Children []Child `gorm:"foreignKey:ParentID"`
}
