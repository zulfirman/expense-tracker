package handlers

import (
	"expenses-tracker/middleware"
	"expenses-tracker/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var categories []models.M_category
	if err := h.db.Where("user_id = ?", userID).Order("name ASC").Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch categories"})
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Generate slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// Check if category with same slug already exists for this user
	var existing models.M_category
	if err := h.db.Where("user_id = ? AND slug = ?", userID, slug).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Category already exists"})
	}

	category := models.M_category{
		UserID:   userID,
		Name:     req.Name,
		Slug:     slug,
		IsActive: true,
	}

	if err := h.db.Create(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create category"})
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var category models.M_category
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Category not found"})
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}
	category.IsActive = req.IsActive

	if err := h.db.Save(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update category"})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id := c.Param("id")

	var category models.M_category
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Category not found"})
	}

	if err := h.db.Delete(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete category"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
