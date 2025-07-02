
package models

import "gorm.io/gorm"

type Reservation struct {
    gorm.Model
    ChildID uint
    SlotID  uint
    Status  string // "pending", "approved", "rejected"
}
