package model

import (
	"time"

	"gorm.io/gorm"
)

type M_expense_template struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"index"`
	Name       string         `json:"name" gorm:"not null"`
	Categories []M_category   `json:"categories" gorm:"many2many:m_expense_template_categories;constraint:OnDelete:CASCADE"`
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Notes      string         `json:"notes" gorm:"type:text"`
	IsActive   bool           `json:"isActive" gorm:"default:true"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
