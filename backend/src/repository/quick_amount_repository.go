package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type QuickAmountRepository struct {
	db *gorm.DB
}

func NewQuickAmountRepository(db *gorm.DB) *QuickAmountRepository {
	return &QuickAmountRepository{db: db}
}

func (r *QuickAmountRepository) GetByUser(userID uint) ([]model.M_quick_amount, error) {
	var list []model.M_quick_amount
	if err := r.db.Where("user_id = ?", userID).Order("value asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *QuickAmountRepository) ReplaceForUser(userID uint, amounts []float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.M_quick_amount{}).Error; err != nil {
			return err
		}
		for _, v := range amounts {
			qa := model.M_quick_amount{
				UserID: userID,
				Value:  v,
			}
			if err := tx.Create(&qa).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
