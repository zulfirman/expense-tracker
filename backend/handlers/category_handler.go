package handlers

import (
	"expenses-tracker/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	var categories []models.Category
	if err := h.db.Where("is_active = ?", true).Order("name ASC").Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch categories"})
	}

	return c.JSON(http.StatusOK, categories)
}

