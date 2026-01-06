package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) GetByUser(userID uint) ([]model.M_expense_template, error) {
	var templates []model.M_expense_template
	if err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *TemplateRepository) Create(t *model.M_expense_template) error {
	return r.db.Create(t).Error
}

func (r *TemplateRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.M_expense_template{}).Error
}


