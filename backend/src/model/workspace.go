package model

import (
	"time"

	"gorm.io/gorm"
)

// M_workspace represents a logical workspace (sheet/project) for a user.
// All financial data (expenses, income, categories, budgets, etc.) are scoped by workspace.
type M_workspace struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"userId" gorm:"index;not null;constraint:OnDelete:CASCADE"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Slug        string         `json:"slug" gorm:"index"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}




