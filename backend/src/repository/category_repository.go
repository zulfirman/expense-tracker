package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll(userID uint, typeFilter string) ([]model.M_category, error) {
	var categories []model.M_category
	query := r.db.Where("user_id = ?", userID)

	// Filter by type if provided
	if typeFilter == "income" || typeFilter == "expense" {
		query = query.Where("type = ?", typeFilter)
	}

	err := query.Order("type ASC, sequence ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetByID(userID uint, id uint) (*model.M_category, error) {
	var category model.M_category
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error
	return &category, err
}

func (r *CategoryRepository) GetByIDs(userID uint, ids []uint) ([]model.M_category, error) {
	var categories []model.M_category
	err := r.db.Where("id IN ? AND user_id = ?", ids, userID).Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetBySlug(userID uint, slug string) (*model.M_category, error) {
	var category model.M_category
	err := r.db.Where("slug = ? AND user_id = ?", slug, userID).First(&category).Error
	return &category, err
}

func (r *CategoryRepository) SlugExists(userID uint, slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.M_category{}).Where("slug = ? AND user_id = ?", slug, userID)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) Create(category *model.M_category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(category *model.M_category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(userID uint, id uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.M_category{}).Error
}
