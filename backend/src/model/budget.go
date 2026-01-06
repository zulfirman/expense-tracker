package model

import (
	"time"

	"gorm.io/gorm"
)

// R_budget represents a per-category monthly budget
type R_budget struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"index"`
	CategoryID uint           `json:"categoryId" gorm:"index"`
	Category   M_category     `json:"category"`
	Month      string         `json:"month" gorm:"type:varchar(7);index"` // YYYY-MM
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
