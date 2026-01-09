package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(t *model.M_refresh_token) error {
	return r.db.Create(t).Error
}

func (r *RefreshTokenRepository) GetByToken(token string) (*model.M_refresh_token, error) {
	var rt model.M_refresh_token
	if err := r.db.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Delete(t *model.M_refresh_token) error {
	return r.db.Delete(t).Error
}

func (r *RefreshTokenRepository) Save(t *model.M_refresh_token) error {
	return r.db.Save(t).Error
}
