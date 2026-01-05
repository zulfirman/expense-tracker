package models

import (
	"time"

	"gorm.io/gorm"
)

type T_income struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"default:null;index"`
	Date      time.Time      `json:"date" gorm:"type:date"`
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2)"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (T_income) TableName() string {
	return "t_incomes"
}

// Balance represents the cumulative income balance
type Balance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"userId" gorm:"default:null;index"`
	Amount    float64        `json:"amount" gorm:"type:decimal(15,2);default:0"`
	Notes     string         `json:"notes" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Balance) TableName() string {
	return "r_balances"
}
