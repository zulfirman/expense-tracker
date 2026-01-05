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
	Name        string  `json:"name"`
	CategoryIDs []uint  `json:"categoryIds"`
	Amount      float64 `json:"amount"`
	Notes       string  `json:"notes"`
}

func (h *TemplateHandler) GetTemplates(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var templates []models.M_expense_template
	if err := h.db.Preload("Categories").Where("user_id = ? AND is_active = ?", userID, true).Order("name ASC").Find(&templates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch templates"})
	}

	// Enrich with category names
	type TemplateResponse struct {
		ID          uint     `json:"id"`
		UserID      uint     `json:"userId"`
		Name        string   `json:"name"`
		CategoryIDs []uint   `json:"categoryIds"`
		Categories  []string `json:"categories"`
		Amount      float64  `json:"amount"`
		Notes       string   `json:"notes"`
		IsActive    bool     `json:"isActive"`
		CreatedAt   string   `json:"createdAt"`
		UpdatedAt   string   `json:"updatedAt"`
	}

	responses := make([]TemplateResponse, len(templates))
	for i, template := range templates {
		categoryIDs := make([]uint, len(template.Categories))
		categoryNames := make([]string, len(template.Categories))
		for j, cat := range template.Categories {
			categoryIDs[j] = cat.ID
			categoryNames[j] = cat.Name
		}
		responses[i] = TemplateResponse{
			ID:          template.ID,
			UserID:      template.UserID,
			Name:        template.Name,
			CategoryIDs: categoryIDs,
			Categories:  categoryNames,
			Amount:      template.Amount,
			Notes:       template.Notes,
			IsActive:    template.IsActive,
			CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *TemplateHandler) CreateTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Load categories
	var categories []models.M_category
	if len(req.CategoryIDs) > 0 {
		if err := h.db.Where("id IN ? AND user_id = ?", req.CategoryIDs, userID).Find(&categories).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to validate categories"})
		}
		if len(categories) != len(req.CategoryIDs) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "One or more category IDs are invalid"})
		}
	}

	template := models.M_expense_template{
		UserID:     userID,
		Name:       req.Name,
		Categories: categories,
		Amount:     req.Amount,
		Notes:      req.Notes,
		IsActive:   true,
	}

	if err := h.db.Create(&template).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create template"})
	}

	// Preload categories for response
	h.db.Preload("Categories").First(&template, template.ID)
	
	// Build response
	categoryIDs := make([]uint, len(template.Categories))
	categoryNames := make([]string, len(template.Categories))
	for i, cat := range template.Categories {
		categoryIDs[i] = cat.ID
		categoryNames[i] = cat.Name
	}
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":          template.ID,
		"userId":      template.UserID,
		"name":        template.Name,
		"categoryIds": categoryIDs,
		"categories":  categoryNames,
		"amount":      template.Amount,
		"notes":       template.Notes,
		"isActive":    template.IsActive,
		"createdAt":   template.CreatedAt,
		"updatedAt":   template.UpdatedAt,
	})
}

func (h *TemplateHandler) DeleteTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid template ID"})
	}

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.M_expense_template{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete template"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}
