package handler

import (
	"net/http"
	"strconv"
	"time"

	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
)

type IncomeHandler struct {
	incomeRepo  *repository.IncomeRepository
	categoryRepo *repository.CategoryRepository
}

func NewIncomeHandler(incomeRepo *repository.IncomeRepository, categoryRepo *repository.CategoryRepository) *IncomeHandler {
	return &IncomeHandler{
		incomeRepo:  incomeRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *IncomeHandler) CreateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		CategoryIDs []uint  `json:"categoryIds"`
		Date        string  `json:"date"`
		Notes       string  `json:"notes"`
		Amount      float64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date"})
	}

	if len(req.CategoryIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "At least one category is required"})
	}

	cats := make([]model.M_category, len(req.CategoryIDs))
	for i, id := range req.CategoryIDs {
		cats[i] = model.M_category{ID: id}
	}

	in := &model.T_income{
		UserID:     cc.UserID,
		Categories: cats,
		Date:       d,
		Notes:      req.Notes,
		Amount:     req.Amount,
	}

	if err := h.incomeRepo.Create(in); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create income"})
	}
	return c.JSON(http.StatusCreated, in)
}

func (h *IncomeHandler) GetBalance(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	b, err := h.incomeRepo.GetBalance(cc.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get balance"})
	}
	if b == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"amount": 0})
	}
	return c.JSON(http.StatusOK, b)
}

func (h *IncomeHandler) UpdateBalance(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		Amount float64 `json:"amount"`
		Notes  string  `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	b := &model.R_balance{
		UserID: cc.UserID,
		Amount: req.Amount,
		Notes:  req.Notes,
	}

	if err := h.incomeRepo.UpsertBalance(b); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update balance"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Balance updated"})
}

func (h *IncomeHandler) GetDateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	dateStr := c.Param("date")
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date"})
	}

	items, err := h.incomeRepo.GetByDate(cc.UserID, d)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch income"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *IncomeHandler) UpdateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	var req struct {
		CategoryIDs []uint  `json:"categoryIds"`
		Date        *string `json:"date"`
		Notes       *string `json:"notes"`
		Amount      *float64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	in, err := h.incomeRepo.GetByID(uint(id), cc.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Income not found"})
	}

	if req.CategoryIDs != nil && len(req.CategoryIDs) > 0 {
		cats := make([]model.M_category, len(req.CategoryIDs))
		for i, catID := range req.CategoryIDs {
			cats[i] = model.M_category{ID: catID}
		}
		if err := h.incomeRepo.ReplaceCategories(in, cats); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update categories"})
		}
	}

	if req.Date != nil {
		if d, err := time.Parse("2006-01-02", *req.Date); err == nil {
			in.Date = d
		}
	}
	if req.Notes != nil {
		in.Notes = *req.Notes
	}
	if req.Amount != nil {
		in.Amount = *req.Amount
	}

	if err := h.incomeRepo.Update(in); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update income"})
	}
	return c.JSON(http.StatusOK, in)
}

func (h *IncomeHandler) DeleteIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	if err := h.incomeRepo.Delete(uint(id), cc.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete income"})
	}
	return c.NoContent(http.StatusNoContent)
}
