package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// StringArray is a custom type for PostgreSQL text array
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return json.Unmarshal([]byte(value.(string)), a)
	}

	return json.Unmarshal(bytes, a)
}

type Expense struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"userId" gorm:"default:null;index"`
	Categories StringArray    `json:"categories" gorm:"type:jsonb"`
	Date       time.Time      `json:"date" gorm:"type:date"`
	Notes      string         `json:"notes" gorm:"type:text"`
	Amount     float64        `json:"amount" gorm:"type:decimal(15,2)"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Expense) TableName() string {
	return "t_expenses"
}
