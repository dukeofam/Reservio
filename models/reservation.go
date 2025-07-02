package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	ChildID uint   `gorm:"index"`
	SlotID  uint   `gorm:"index"`
	Status  string // "pending", "approved", "rejected"
}
