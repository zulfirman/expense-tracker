package model

import (
	"time"

	"gorm.io/gorm"
)

// R_budget represents a per-category monthly budget
type R_budget struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"index;constraint:OnDelete:CASCADE"`
	CategoryID uint           `json:"categoryId" gorm:"index;constraint:OnDelete:CASCADE"`
	Category   M_category     `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Month      string         `json:"month" gorm:"type:varchar(7);index"` // YYYY-MM
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
