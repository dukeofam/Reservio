package models

import "gorm.io/gorm"

// PasswordResetToken represents a single-use, time-limited token for resetting a user password.
// Tokens are deleted after use or when expired to prevent reuse.
type PasswordResetToken struct {
	gorm.Model
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex;size:191"` // size 191 works with most MySQL indexes; fine for Postgres too
	ExpiresAt int64  `gorm:"index"`                // Unix timestamp (seconds)
}
