package models

import "gorm.io/gorm"

type Child struct {
	gorm.Model
	Name      string
	Age       int
	ParentID  uint   `gorm:"index"`
	Birthdate string `json:"birthdate"`
	Parent    *User  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}
