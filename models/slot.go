
package models

import "gorm.io/gorm"

type Slot struct {
    gorm.Model
    Date     string `gorm:"unique"`
    Capacity int
}
