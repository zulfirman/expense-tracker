package handlers

import (
	"expenses-tracker/middleware"
	"expenses-tracker/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BudgetHandler struct {
	db *gorm.DB
}

func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

type CreateBudgetRequest struct {
	CategoryID uint    `json:"categoryId" validate:"required"`
	Month      string  `json:"month" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

func (h *BudgetHandler) GetBudgets(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	month := c.QueryParam("month")

	if month == "" {
		// Default to current month
		month = time.Now().Format("2006-01")
	}

	var budgets []models.R_budget
	query := h.db.Where("user_id = ? AND month = ?", userID, month)
	if err := query.Find(&budgets).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch budgets"})
	}

	// Get categories for enrichment
	var categoryIDs []uint
	for _, budget := range budgets {
		categoryIDs = append(categoryIDs, budget.CategoryID)
	}
	var categories []models.M_category
	if len(categoryIDs) > 0 {
		h.db.Where("id IN ? AND user_id = ?", categoryIDs, userID).Find(&categories)
	}
	categoryMap := make(map[uint]models.M_category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	// Enrich budgets with category info
	type BudgetResponse struct {
		ID           uint    `json:"id"`
		UserID       uint    `json:"userId"`
		CategoryID   uint    `json:"categoryId"`
		CategoryName string  `json:"categoryName"`
		CategorySlug string  `json:"categorySlug"`
		Month        string  `json:"month"`
		Amount       float64 `json:"amount"`
		CreatedAt    string  `json:"createdAt"`
		UpdatedAt    string  `json:"updatedAt"`
	}

	responses := make([]BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = BudgetResponse{
			ID:           budget.ID,
			UserID:       budget.UserID,
			CategoryID:   budget.CategoryID,
			CategoryName: categoryMap[budget.CategoryID].Name,
			CategorySlug: categoryMap[budget.CategoryID].Slug,
			Month:        budget.Month,
			Amount:       budget.Amount,
			CreatedAt:    budget.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    budget.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateBudgetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Validate category exists
	var category models.M_category
	if err := h.db.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Category not found"})
	}

	// Check if budget already exists for this category and month
	var existing models.R_budget
	if err := h.db.Where("user_id = ? AND category_id = ? AND month = ?", userID, req.CategoryID, req.Month).First(&existing).Error; err == nil {
		// Update existing budget
		existing.Amount = req.Amount
		if err := h.db.Save(&existing).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update budget"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"id":           existing.ID,
			"userId":       existing.UserID,
			"categoryId":   existing.CategoryID,
			"categoryName": category.Name,
			"categorySlug": category.Slug,
			"month":        existing.Month,
			"amount":       existing.Amount,
			"createdAt":    existing.CreatedAt,
			"updatedAt":    existing.UpdatedAt,
		})
	}

	// Create new budget
	budget := models.R_budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Month:      req.Month,
		Amount:     req.Amount,
	}

	if err := h.db.Create(&budget).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create budget"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":           budget.ID,
		"userId":       budget.UserID,
		"categoryId":   budget.CategoryID,
		"categoryName": category.Name,
		"categorySlug": category.Slug,
		"month":        budget.Month,
		"amount":       budget.Amount,
		"createdAt":    budget.CreatedAt,
		"updatedAt":    budget.UpdatedAt,
	})
}

func (h *BudgetHandler) DeleteBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	categoryIDStr := c.Param("categoryId")
	month := c.QueryParam("month")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	var budget models.R_budget
	if err := h.db.Where("user_id = ? AND category_id = ? AND month = ?", userID, categoryID, month).First(&budget).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Budget not found"})
	}

	if err := h.db.Delete(&budget).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete budget"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Budget deleted successfully"})
}
