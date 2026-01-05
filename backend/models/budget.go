package models

import (
	"time"

	"gorm.io/gorm"
)

type R_budget struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"default:null;index"`
	CategoryID uint           `json:"categoryId" gorm:"not null;index"`
	Category   M_category     `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	Month      string         `json:"month" gorm:"type:varchar(7);not null;index"` // Format: YYYY-MM
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (R_budget) TableName() string {
	return "r_budgets"
}
