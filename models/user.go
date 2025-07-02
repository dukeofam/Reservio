package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string  // "parent" or "admin"
	Children []Child `gorm:"foreignKey:ParentID"`
}
