package model

type M_quick_amount struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	UserID uint    `json:"userId" gorm:"index;constraint:OnDelete:CASCADE"`
	WorkspaceID uint `json:"workspaceId" gorm:"index;not null;default:0"`
	Value  float64 `json:"value"`
}
