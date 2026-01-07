package repository

import (
	"expenses-tracker/src/model"
	"expenses-tracker/src/utils"

	"gorm.io/gorm"
)

// SeedDefaultCategories creates default income and expense categories for a new user
func SeedDefaultCategories(db *gorm.DB, userID uint) error {
	defaultExpenseCategories := []string{
		"Food & Dining",
		"Transportation",
		"Shopping",
		"Bills & Utilities",
		"Entertainment",
		"Subscriptions",
		"Other Expenses",
	}

	defaultIncomeCategories := []string{
		"Salary",
		"Freelance",
		"Investment",
		"Other Income",
	}

	categories := make([]model.M_category, 0, len(defaultExpenseCategories)+len(defaultIncomeCategories))

	// Add expense categories
	for _, name := range defaultExpenseCategories {
		slug := utils.GenerateSlug(name)
		categories = append(categories, model.M_category{
			UserID:   userID,
			Name:     name,
			Slug:     slug,
			Type:     "expense",
			IsActive: true,
		})
	}

	// Add income categories
	for _, name := range defaultIncomeCategories {
		slug := utils.GenerateSlug(name)
		categories = append(categories, model.M_category{
			UserID:   userID,
			Name:     name,
			Slug:     slug,
			Type:     "income",
			IsActive: true,
		})
	}

	// Batch insert all categories
	return db.Create(&categories).Error
}
