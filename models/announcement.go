package models

import "gorm.io/gorm"

// Announcement represents a public announcement posted by an admin.
// CreatedAt is used as the published date.
// AuthorID references the User (admin) who created it.
// If Author record is deleted, announcements remain with NULL author.
type Announcement struct {
	gorm.Model
	Title    string `gorm:"size:200" json:"title"`
	Content  string `gorm:"type:text" json:"content"`
	AuthorID uint   `json:"author_id"`
	Author   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author,omitempty"`
}
