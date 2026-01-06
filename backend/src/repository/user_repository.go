package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.M_user) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*model.M_user, error) {
	var user model.M_user
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*model.M_user, error) {
	var user model.M_user
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.M_user) error {
	return r.db.Save(user).Error
}

// EmailExists returns true if there is another user (not excludedID) with the given email.
func (r *UserRepository) EmailExists(email string, excludedID uint) (bool, error) {
	var count int64
	if err := r.db.
		Model(&model.M_user{}).
		Where("email = ? AND id <> ?", email, excludedID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
