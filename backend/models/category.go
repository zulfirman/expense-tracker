package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Slug      string         `json:"slug" gorm:"unique;not null"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	UserID    uint           `json:"userId" gorm:"default:null;index"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Category) TableName() string {
	return "m_categories"
}
