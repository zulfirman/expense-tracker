package model

import (
	"time"

	"gorm.io/gorm"
)

type T_expense struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"index"`
	Categories []M_category   `json:"categories" gorm:"many2many:t_expense_categories;constraint:OnDelete:CASCADE"`
	Date       time.Time      `json:"date" gorm:"type:date"`
	Notes      string         `json:"notes" gorm:"type:text"`
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
