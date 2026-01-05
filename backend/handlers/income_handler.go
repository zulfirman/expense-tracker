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

type IncomeHandler struct {
	db *gorm.DB
}

func NewIncomeHandler(db *gorm.DB) *IncomeHandler {
	return &IncomeHandler{db: db}
}

type CreateIncomeRequest struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Notes  string  `json:"notes"`
}

type UpdateBalanceRequest struct {
	Amount float64 `json:"amount"`
	Notes  string  `json:"notes"`
}

func (h *IncomeHandler) CreateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateIncomeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
	}

	income := models.Income{
		UserID: userID,
		Date:   date,
		Amount: req.Amount,
		Notes:  req.Notes,
	}

	if err := h.db.Create(&income).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create income"})
	}

	// Update balance
	var balance models.Balance
	if err := h.db.FirstOrCreate(&balance, models.Balance{UserID: userID}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update balance"})
	}
	balance.Amount += req.Amount
	if err := h.db.Save(&balance).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update balance"})
	}

	return c.JSON(http.StatusCreated, income)
}

func (h *IncomeHandler) GetBalance(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var balance models.Balance
	if err := h.db.FirstOrCreate(&balance, models.Balance{UserID: userID}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch balance"})
	}

	return c.JSON(http.StatusOK, balance)
}

func (h *IncomeHandler) UpdateBalance(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req UpdateBalanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var balance models.Balance
	if err := h.db.FirstOrCreate(&balance, models.Balance{UserID: userID}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch balance"})
	}

	balance.Amount = req.Amount
	balance.Notes = req.Notes

	if err := h.db.Save(&balance).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update balance"})
	}

	return c.JSON(http.StatusOK, balance)
}

func (h *IncomeHandler) GetDateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
	}

	var incomes []models.Income
	err = h.db.Where("user_id = ? AND date = ?", userID, date).Order("created_at DESC").Find(&incomes).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch income"})
	}

	return c.JSON(http.StatusOK, incomes)
}

func (h *IncomeHandler) UpdateIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid income ID"})
	}

	var req CreateIncomeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var income models.Income
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&income).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Income not found"})
	}

	// Calculate difference for balance update
	oldAmount := income.Amount
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
	}

	income.Date = date
	income.Amount = req.Amount
	income.Notes = req.Notes

	if err := h.db.Save(&income).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update income"})
	}

	// Update balance with difference
	var balance models.Balance
	if err := h.db.FirstOrCreate(&balance, models.Balance{UserID: userID}).Error; err == nil {
		balance.Amount = balance.Amount - oldAmount + req.Amount
		h.db.Save(&balance)
	}

	return c.JSON(http.StatusOK, income)
}

func (h *IncomeHandler) DeleteIncome(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid income ID"})
	}

	var income models.Income
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&income).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Income not found"})
	}

	amount := income.Amount
	if err := h.db.Unscoped().Delete(&income).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete income"})
	}

	// Update balance
	var balance models.Balance
	if err := h.db.FirstOrCreate(&balance, models.Balance{UserID: userID}).Error; err == nil {
		balance.Amount -= amount
		h.db.Save(&balance)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Income deleted successfully"})
}

