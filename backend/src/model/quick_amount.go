package model

type M_quick_amount struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	UserID uint    `json:"userId" gorm:"index"`
	Value  float64 `json:"value"`
}
