package handler

import (
	"net/http"
	"strconv"

	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	budgetRepo   *repository.BudgetRepository
	categoryRepo *repository.CategoryRepository
}

func NewBudgetHandler(budgetRepo *repository.BudgetRepository, categoryRepo *repository.CategoryRepository) *BudgetHandler {
	return &BudgetHandler{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *BudgetHandler) GetBudgets(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	month := c.QueryParam("month")

	var items []model.R_budget
	var err error
	if month != "" {
		items, err = h.budgetRepo.GetByUserAndMonth(cc.UserID, month)
	} else {
		items, err = h.budgetRepo.GetByUser(cc.UserID)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch budgets"})
	}

	// Shape response to what frontend expects:
	// [{ categoryId, categoryName, amount, month }]
	type resp struct {
		CategoryID   uint    `json:"categoryId"`
		CategoryName string  `json:"categoryName"`
		Amount       float64 `json:"amount"`
		Month        string  `json:"month"`
	}
	out := make([]resp, 0, len(items))
	for _, b := range items {
		name := ""
		if b.Category.ID != 0 {
			name = b.Category.Name
		}
		out = append(out, resp{
			CategoryID:   b.CategoryID,
			CategoryName: name,
			Amount:       b.Amount,
			Month:        b.Month,
		})
	}

	return c.JSON(http.StatusOK, out)
}

func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		CategoryID uint    `json:"categoryId"`
		Month      string  `json:"month"`
		Amount     float64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	b := &model.R_budget{
		UserID:     cc.UserID,
		CategoryID: req.CategoryID,
		Month:      req.Month,
		Amount:     req.Amount,
	}

	if err := h.budgetRepo.Create(b); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create budget"})
	}
	return c.JSON(http.StatusCreated, b)
}

// CopyBudgets copies all budgets from sourceMonth to targetMonth for the current user.
func (h *BudgetHandler) CopyBudgets(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		FromMonth string `json:"fromMonth"`
		ToMonth   string `json:"toMonth"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Simple implementation: read all source, create new rows for target
	items, err := h.budgetRepo.GetByUserAndMonth(cc.UserID, req.FromMonth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch budgets"})
	}

	for _, b := range items {
		nb := model.R_budget{
			UserID:     cc.UserID,
			CategoryID: b.CategoryID,
			Month:      req.ToMonth,
			Amount:     b.Amount,
		}
		if err := h.budgetRepo.Create(&nb); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to copy budget"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Budgets copied"})
}

func (h *BudgetHandler) GetLatestBudgetMonth(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	month, err := h.budgetRepo.LatestMonth(cc.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get latest budget month"})
	}
	return c.JSON(http.StatusOK, map[string]string{"month": month})
}

func (h *BudgetHandler) DeleteBudget(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("categoryId")
	catID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	month := c.QueryParam("month")

	if err := h.budgetRepo.DeleteByCategory(cc.UserID, uint(catID), month); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete budget"})
	}
	return c.NoContent(http.StatusNoContent)
}


