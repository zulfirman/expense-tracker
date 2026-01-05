package handlers

import (
	"expenses-tracker/middleware"
	"expenses-tracker/models"
	"net/http"
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
	CategorySlug string  `json:"categorySlug" validate:"required"`
	Month        string  `json:"month" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
}

func (h *BudgetHandler) GetBudgets(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	month := c.QueryParam("month")

	if month == "" {
		// Default to current month
		month = time.Now().Format("2006-01")
	}

	var budgets []models.Budget
	query := h.db.Where("user_id = ? AND month = ?", userID, month)
	if err := query.Find(&budgets).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch budgets"})
	}

	return c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateBudgetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Check if budget already exists for this category and month
	var existing models.Budget
	if err := h.db.Where("user_id = ? AND category_slug = ? AND month = ?", userID, req.CategorySlug, req.Month).First(&existing).Error; err == nil {
		// Update existing budget
		existing.Amount = req.Amount
		if err := h.db.Save(&existing).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update budget"})
		}
		return c.JSON(http.StatusOK, existing)
	}

	// Create new budget
	budget := models.Budget{
		UserID:       userID,
		CategorySlug: req.CategorySlug,
		Month:        req.Month,
		Amount:       req.Amount,
	}

	if err := h.db.Create(&budget).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create budget"})
	}

	return c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) DeleteBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	categorySlug := c.Param("categorySlug")
	month := c.QueryParam("month")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var budget models.Budget
	if err := h.db.Where("user_id = ? AND category_slug = ? AND month = ?", userID, categorySlug, month).First(&budget).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Budget not found"})
	}

	if err := h.db.Delete(&budget).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete budget"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Budget deleted successfully"})
}

