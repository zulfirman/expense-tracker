package model

import (
	"time"

	"gorm.io/gorm"
)

type M_refresh_token struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"` // refresh token string (can be hashed or raw)
	DeviceID  string         `json:"deviceId" gorm:"size:255"`
	UsedCount int            `json:"usedCount" gorm:"default:0"` // how many times this refresh token has been used
	ExpiresAt time.Time      `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
