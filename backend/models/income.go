package models

import (
	"time"

	"gorm.io/gorm"
)

type Income struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Date      time.Time      `json:"date" gorm:"type:date"`
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Balance represents the cumulative income balance
type Balance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2);default:0"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

