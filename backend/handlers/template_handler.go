package handlers

import (
	"expenses-tracker/middleware"
	"expenses-tracker/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TemplateHandler struct {
	db *gorm.DB
}

func NewTemplateHandler(db *gorm.DB) *TemplateHandler {
	return &TemplateHandler{db: db}
}

type CreateTemplateRequest struct {
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Amount     float64  `json:"amount"`
	Notes      string   `json:"notes"`
}

func (h *TemplateHandler) GetTemplates(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var templates []models.ExpenseTemplate
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).Order("name ASC").Find(&templates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch templates"})
	}

	return c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) CreateTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	template := models.ExpenseTemplate{
		UserID:     userID,
		Name:       req.Name,
		Categories: req.Categories,
		Amount:     req.Amount,
		Notes:      req.Notes,
		IsActive:   true,
	}

	if err := h.db.Create(&template).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create template"})
	}

	return c.JSON(http.StatusCreated, template)
}

func (h *TemplateHandler) DeleteTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid template ID"})
	}

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.ExpenseTemplate{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete template"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}

