package repository

import (
	"expenses-tracker/src/model"

	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ws *model.M_workspace) error {
	return r.db.Create(ws).Error
}

func (r *WorkspaceRepository) GetByID(userID uint, id uint) (*model.M_workspace, error) {
	var ws model.M_workspace
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&ws).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}

func (r *WorkspaceRepository) ListByUser(userID uint, q string) ([]model.M_workspace, error) {
	var list []model.M_workspace
	query := r.db.Where("user_id = ?", userID)
	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}
	err := query.Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *WorkspaceRepository) Update(ws *model.M_workspace) error {
	return r.db.Save(ws).Error
}





