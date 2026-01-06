package model

import (
	"time"

	"gorm.io/gorm"
)

type M_category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Slug      string         `json:"slug" gorm:"default:null;index:idx_category_slug_user"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	UserID    uint           `json:"userId" gorm:"default:null;index:idx_category_user_id;index:idx_category_slug_user"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
