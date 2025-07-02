package models

import "gorm.io/gorm"

type Child struct {
	gorm.Model
	Name     string
	Age      int
	ParentID uint `gorm:"index"`
}
