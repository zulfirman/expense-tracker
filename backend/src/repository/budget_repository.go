package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) GetByUser(userID uint) ([]model.R_budget, error) {
	var budgets []model.R_budget
	if err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("month DESC, category_id ASC").
		Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) GetByUserAndMonth(userID uint, month string) ([]model.R_budget, error) {
	var budgets []model.R_budget
	if err := r.db.Preload("Category").
		Where("user_id = ? AND month = ?", userID, month).
		Order("category_id ASC").
		Find(&budgets).Error; err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *BudgetRepository) Create(budget *model.R_budget) error {
	return r.db.Create(budget).Error
}

func (r *BudgetRepository) DeleteByCategory(userID, categoryID uint, month string) error {
	db := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID)
	if month != "" {
		db = db.Where("month = ?", month)
	}
	return db.
		Delete(&model.R_budget{}).Error
}

func (r *BudgetRepository) LatestMonth(userID uint) (string, error) {
	type Row struct {
		Month string
	}
	var row Row
	if err := r.db.
		Model(&model.R_budget{}).
		Select("month").
		Where("user_id = ?", userID).
		Order("month DESC").
		Limit(1).
		Scan(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return row.Month, nil
}


