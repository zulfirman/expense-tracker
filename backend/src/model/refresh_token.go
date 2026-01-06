package model

import (
	"time"

	"gorm.io/gorm"
)

type M_refresh_token struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
