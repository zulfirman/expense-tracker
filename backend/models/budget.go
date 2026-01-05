package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"userId" gorm:"default:null;index"`
	CategorySlug string         `json:"categorySlug" gorm:"not null"`
	Month        string         `json:"month" gorm:"type:varchar(7);not null;index"` // Format: YYYY-MM
	Amount       float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Budget) TableName() string {
	return "r_budgets"
}
