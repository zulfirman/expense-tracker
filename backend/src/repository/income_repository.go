package repository

import (
	"errors"
	"time"

	"expenses-tracker/src/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IncomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) GetByID(id, userID uint) (*model.T_income, error) {
	var in model.T_income
	if err := r.db.Preload("Categories").Where("id = ? AND user_id = ?", id, userID).First(&in).Error; err != nil {
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
	if err := r.db.Preload("Categories").
		Where("user_id = ? AND date = ?", userID, date).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *IncomeRepository) ReplaceCategories(income *model.T_income, categories []model.M_category) error {
	return r.db.Model(income).Association("Categories").Replace(categories)
}

func (r *IncomeRepository) GetBalance(userID uint) (*model.R_balance, error) {
	var b model.R_balance

	// Always derive the latest balance from income and expenses, and upsert it atomically.
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var incomeTotal float64
		var expenseTotal float64

		if err := tx.
			Model(&model.T_income{}).
			Where("user_id = ?", userID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&incomeTotal).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&model.T_expense{}).
			Where("user_id = ?", userID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&expenseTotal).Error; err != nil {
			return err
		}

		b = model.R_balance{
			UserID: userID,
			Amount: incomeTotal - expenseTotal,
		}

		// Upsert by user_id (requires unique index on user_id).
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount", "updated_at"}),
		}).Create(&b).Error
	}); err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *IncomeRepository) UpsertBalance(b *model.R_balance) error {
	// upsert based on user_id with row-level locking to avoid races
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.R_balance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", b.UserID).
			First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(b).Error
			}
			return err
		}
		existing.Amount = b.Amount
		existing.Notes = b.Notes
		return tx.Save(&existing).Error
	})
}
