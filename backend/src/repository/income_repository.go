package repository

import (
	"errors"
	"time"

	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type IncomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) GetByID(id, userID uint) (*model.T_income, error) {
	var in model.T_income
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&in).Error; err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *IncomeRepository) Create(income *model.T_income) error {
	return r.db.Create(income).Error
}

func (r *IncomeRepository) Update(income *model.T_income) error {
	return r.db.Save(income).Error
}

func (r *IncomeRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.T_income{}).Error
}

func (r *IncomeRepository) GetByDate(userID uint, date time.Time) ([]model.T_income, error) {
	var items []model.T_income
	if err := r.db.
		Where("user_id = ? AND date = ?", userID, date).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *IncomeRepository) GetBalance(userID uint) (*model.R_balance, error) {
	var b model.R_balance
	if err := r.db.Where("user_id = ?", userID).First(&b).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// No stored balance yet: derive it from total income minus total expenses
		var incomeTotal float64
		var expenseTotal float64

		if err := r.db.
			Model(&model.T_income{}).
			Where("user_id = ?", userID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&incomeTotal).Error; err != nil {
			return nil, err
		}

		if err := r.db.
			Model(&model.T_expense{}).
			Where("user_id = ?", userID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&expenseTotal).Error; err != nil {
			return nil, err
		}

		// Return a synthetic balance object without persisting it
		b = model.R_balance{
			UserID: userID,
			Amount: incomeTotal - expenseTotal,
		}
		return &b, nil
	}
	return &b, nil
}

func (r *IncomeRepository) UpsertBalance(b *model.R_balance) error {
	// simple upsert based on user_id
	var existing model.R_balance
	if err := r.db.Where("user_id = ?", b.UserID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(b).Error
		}
		return err
	}
	existing.Amount = b.Amount
	existing.Notes = b.Notes
	return r.db.Save(&existing).Error
}
