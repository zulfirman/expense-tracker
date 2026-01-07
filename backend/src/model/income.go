package model

import (
	"time"

	"gorm.io/gorm"
)

type T_income struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"index;constraint:OnDelete:CASCADE"`
	Categories []M_category   `json:"categories" gorm:"many2many:t_income_categories;constraint:OnDelete:CASCADE"`
	Date       time.Time      `json:"date" gorm:"type:date"`
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Notes      string         `json:"notes" gorm:"type:text"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type R_balance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"uniqueIndex;constraint:OnDelete:CASCADE"` // one balance row per user
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2);default:0"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
