package repository

import (
	"time"

	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(expense *model.T_expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseRepository) Update(expense *model.T_expense) error {
	return r.db.Save(expense).Error
}

func (r *ExpenseRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.T_expense{}).Error
}

func (r *ExpenseRepository) GetByID(id uint, userID uint) (*model.T_expense, error) {
	var e model.T_expense
	if err := r.db.Preload("Categories").Where("id = ? AND user_id = ?", id, userID).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// GetMonths returns distinct YYYY-MM strings where the user has expenses.
func (r *ExpenseRepository) GetMonths(userID uint) ([]string, error) {
	type Row struct {
		Month string
	}
	var rows []Row
	if err := r.db.
		Model(&model.T_expense{}).
		Select("to_char(date, 'YYYY-MM') as month").
		Where("user_id = ?", userID).
		Group("month").
		Order("month DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	months := make([]string, 0, len(rows))
	for _, r := range rows {
		months = append(months, r.Month)
	}
	return months, nil
}

// GetByMonth returns all expenses for a given YYYY-MM.
func (r *ExpenseRepository) GetByMonth(userID uint, month string) ([]model.T_expense, error) {
	var expenses []model.T_expense
	if err := r.db.
		Preload("Categories").
		Where("user_id = ? AND to_char(date, 'YYYY-MM') = ?", userID, month).
		Order("date DESC, id DESC").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetByDate returns all expenses for a specific date (YYYY-MM-DD).
func (r *ExpenseRepository) GetByDate(userID uint, date time.Time) ([]model.T_expense, error) {
	var expenses []model.T_expense
	if err := r.db.
		Preload("Categories").
		Where("user_id = ? AND date = ?", userID, date).
		Order("id DESC").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// Search by notes substring with optional filters.
func (r *ExpenseRepository) Search(userID uint, query string, categoryID *uint, dateFrom, dateTo *time.Time) ([]model.T_expense, error) {
	var expenses []model.T_expense

	db := r.db.Preload("Categories").Where("user_id = ?", userID)
	if query != "" {
		db = db.Where("notes ILIKE ?", "%"+query+"%")
	}
	if dateFrom != nil {
		db = db.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		db = db.Where("date <= ?", *dateTo)
	}
	if categoryID != nil {
		// Join through many2many table to filter by category
		db = db.Joins("JOIN t_expense_categories ON t_expense_categories.t_expense_id = t_expenses.id").
			Where("t_expense_categories.m_category_id = ?", *categoryID)
	}

	if err := db.
		Order("date DESC, id DESC").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// ReplaceCategories replaces the categories association for an expense.
func (r *ExpenseRepository) ReplaceCategories(expense *model.T_expense, categories []model.M_category) error {
	return r.db.Model(expense).Association("Categories").Replace(categories)
}


