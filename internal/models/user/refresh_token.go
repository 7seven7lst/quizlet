package user

import (
	"time"
	"gorm.io/gorm"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Revoked   bool           `gorm:"not null;default:false" json:"revoked"`
} 