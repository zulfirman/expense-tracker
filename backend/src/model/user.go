package model

import (
	"time"

	"gorm.io/gorm"
)

type M_user struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	Name                  string         `json:"name" gorm:"not null"`
	Email                 string         `json:"email" gorm:"uniqueIndex;not null"`
	Password              string         `json:"-" gorm:"not null"`
	Currency              string         `json:"currency" gorm:"default:'IDR'"` // IDR, USD, EUR, JPY
	FirstSigninCompleted  bool           `json:"firstSigninCompleted" gorm:"default:false"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}
