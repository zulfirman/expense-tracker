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

type ExpenseHandler struct {
	db *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{db: db}
}

type CreateExpenseRequest struct {
	Categories []string `json:"categories"`
	Date       string   `json:"date"`
	Notes      string   `json:"notes"`
	Amount     float64  `json:"amount"`
}

type UpdateExpenseRequest struct {
	Categories []string `json:"categories"`
	Notes      string   `json:"notes"`
	Amount     float64  `json:"amount"`
}

func (h *ExpenseHandler) CreateExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
	}

	expense := models.Expense{
		UserID:     userID,
		Categories: req.Categories,
		Date:       date,
		Notes:      req.Notes,
		Amount:     req.Amount,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create expense"})
	}

	return c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) GetMonths(c echo.Context) error {
	userID := middleware.GetUserID(c)

	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	threeMonthsAgo := currentMonth.AddDate(0, -3, 0)

	type MonthResult struct {
		Month string  `json:"month"`
		Total float64 `json:"total"`
	}

	var results []MonthResult
	err := h.db.Model(&models.Expense{}).
		Select("TO_CHAR(date, 'YYYY-MM') as month, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND date >= ?", userID, threeMonthsAgo).
		Group("TO_CHAR(date, 'YYYY-MM')").
		Order("month DESC").
		Scan(&results).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch months"})
	}

	// Generate all months from 3 months ago to current month
	allMonths := make(map[string]bool)
	monthMap := make(map[string]float64)
	for _, result := range results {
		allMonths[result.Month] = true
		monthMap[result.Month] = result.Total
	}

	// Add all months in range (including empty ones)
	iterMonth := threeMonthsAgo
	for iterMonth.Before(currentMonth) || iterMonth.Equal(currentMonth) {
		monthKey := iterMonth.Format("2006-01")
		allMonths[monthKey] = true
		if _, exists := monthMap[monthKey]; !exists {
			monthMap[monthKey] = 0
		}
		iterMonth = iterMonth.AddDate(0, 1, 0)
	}

	type MonthWithDates struct {
		Month string      `json:"month"`
		Total float64     `json:"total"`
		Dates []DateTotal `json:"dates"`
	}

	// Build response with all months
	var monthsWithDates []MonthWithDates
	iterMonth = currentMonth
	for iterMonth.After(threeMonthsAgo) || iterMonth.Equal(threeMonthsAgo) {
		monthKey := iterMonth.Format("2006-01")

		type DateResult struct {
			Date  time.Time `json:"date"`
			Total float64   `json:"total"`
		}
		var dateResults []DateResult
		h.db.Model(&models.Expense{}).
			Select("date, COALESCE(SUM(amount), 0) as total").
			Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, monthKey).
			Group("date").
			Order("date DESC").
			Scan(&dateResults)

		// Get all unique dates (from expenses and income) for this month
		dateMap := make(map[time.Time]struct {
			Expense float64
			Income  float64
		})

		// Add expenses to date map
		for _, dr := range dateResults {
			dateMap[dr.Date] = struct {
				Expense float64
				Income  float64
			}{Expense: dr.Total, Income: 0}
		}

		// Get income for this month
		var incomeResults []struct {
			Date  time.Time `json:"date"`
			Total float64   `json:"total"`
		}
		h.db.Model(&models.Income{}).
			Select("date, COALESCE(SUM(amount), 0) as total").
			Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, monthKey).
			Group("date").
			Scan(&incomeResults)

		// Add income to date map
		for _, ir := range incomeResults {
			if existing, ok := dateMap[ir.Date]; ok {
				existing.Income = ir.Total
				dateMap[ir.Date] = existing
			} else {
				dateMap[ir.Date] = struct {
					Expense float64
					Income  float64
				}{Expense: 0, Income: ir.Total}
			}
		}

		// Convert to DateTotal slice with net total (income - expense) and type info
		type DateTotalWithType struct {
			Date       string  `json:"date"`
			Total      float64 `json:"total"`
			HasIncome  bool    `json:"hasIncome"`
			HasExpense bool    `json:"hasExpense"`
		}
		dateTotalsWithType := make([]DateTotalWithType, 0, len(dateMap))
		for date, amounts := range dateMap {
			netTotal := amounts.Income - amounts.Expense
			dateTotalsWithType = append(dateTotalsWithType, DateTotalWithType{
				Date:       date.Format("2006-01-02"),
				Total:      netTotal,
				HasIncome:  amounts.Income > 0,
				HasExpense: amounts.Expense > 0,
			})
		}

		// Sort by date descending (31 to 1)
		for i := 0; i < len(dateTotalsWithType); i++ {
			for j := i + 1; j < len(dateTotalsWithType); j++ {
				if dateTotalsWithType[i].Date < dateTotalsWithType[j].Date {
					dateTotalsWithType[i], dateTotalsWithType[j] = dateTotalsWithType[j], dateTotalsWithType[i]
				}
			}
		}

		// Convert to DateTotal format with income/expense flags
		dateTotals := make([]DateTotal, len(dateTotalsWithType))
		for j, dt := range dateTotalsWithType {
			dateTotals[j] = DateTotal{
				Date:       dt.Date,
				Total:      dt.Total,
				HasIncome:  dt.HasIncome,
				HasExpense: dt.HasExpense,
			}
		}

		// Use expense total only (no income calculation)
		expenseTotal := monthMap[monthKey]

		monthsWithDates = append(monthsWithDates, MonthWithDates{
			Month: monthKey,
			Total: expenseTotal,
			Dates: dateTotals,
		})

		iterMonth = iterMonth.AddDate(0, -1, 0)
		if iterMonth.Before(threeMonthsAgo) {
			break
		}
	}

	return c.JSON(http.StatusOK, monthsWithDates)
}

type DateTotal struct {
	Date       string  `json:"date"`
	Total      float64 `json:"total"`
	HasIncome  bool    `json:"hasIncome,omitempty"`
	HasExpense bool    `json:"hasExpense,omitempty"`
}

func (h *ExpenseHandler) GetMonthDetails(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	month := c.Param("month")

	type CategoryTotal struct {
		Category string  `json:"category"`
		Total    float64 `json:"total"`
	}

	type DailyTotal struct {
		Date     string  `json:"date"`
		Income   float64 `json:"income"`
		Expense  float64 `json:"expense"`
		NetTotal float64 `json:"netTotal"`
	}

	// Get all expenses for the month
	var expenses []models.Expense
	if err := h.db.Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).Find(&expenses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch month details"})
	}

	// Get all income for the month
	var incomes []models.Income
	if err := h.db.Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).Find(&incomes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch income details"})
	}

	// Aggregate by category
	categoryMap := make(map[string]float64)
	dailyMap := make(map[string]struct {
		Income  float64
		Expense float64
	})

	// Process expenses
	for _, expense := range expenses {
		dateKey := expense.Date.Format("2006-01-02")
		if existing, ok := dailyMap[dateKey]; ok {
			existing.Expense += expense.Amount
			dailyMap[dateKey] = existing
		} else {
			dailyMap[dateKey] = struct {
				Income  float64
				Expense float64
			}{Expense: expense.Amount}
		}
		for _, category := range expense.Categories {
			categoryMap[category] += expense.Amount
		}
	}

	// Process income
	for _, income := range incomes {
		dateKey := income.Date.Format("2006-01-02")
		if existing, ok := dailyMap[dateKey]; ok {
			existing.Income += income.Amount
			dailyMap[dateKey] = existing
		} else {
			dailyMap[dateKey] = struct {
				Income  float64
				Expense float64
			}{Income: income.Amount}
		}
	}

	categoryResults := make([]CategoryTotal, 0, len(categoryMap))
	for category, total := range categoryMap {
		categoryResults = append(categoryResults, CategoryTotal{
			Category: category,
			Total:    total,
		})
	}

	dailyResults := make([]DailyTotal, 0, len(dailyMap))
	for date, amounts := range dailyMap {
		dailyResults = append(dailyResults, DailyTotal{
			Date:     date,
			Income:   amounts.Income,
			Expense:  amounts.Expense,
			NetTotal: amounts.Income - amounts.Expense,
		})
	}

	// Sort daily results by date
	for i := 0; i < len(dailyResults)-1; i++ {
		for j := i + 1; j < len(dailyResults); j++ {
			if dailyResults[i].Date > dailyResults[j].Date {
				dailyResults[i], dailyResults[j] = dailyResults[j], dailyResults[i]
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categoryResults,
		"daily":      dailyResults,
	})
}

func (h *ExpenseHandler) GetDateExpenses(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
	}

	var expenses []models.Expense
	err = h.db.Where("user_id = ? AND date = ?", userID, date).Order("created_at DESC").Find(&expenses).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch expenses"})
	}

	return c.JSON(http.StatusOK, expenses)
}

func (h *ExpenseHandler) SearchExpenses(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	query := h.db.Model(&models.Expense{}).Where("user_id = ?", userID)

	// Search by text (notes or categories)
	searchQuery := c.QueryParam("q")
	if searchQuery != "" {
		query = query.Where("notes ILIKE ? OR EXISTS (SELECT 1 FROM jsonb_array_elements_text(categories) AS cat WHERE cat ILIKE ?)", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Filter by category
	category := c.QueryParam("category")
	if category != "" {
		query = query.Where("EXISTS (SELECT 1 FROM jsonb_array_elements_text(categories) AS cat WHERE cat = ?)", category)
	}

	// Filter by date range
	dateFrom := c.QueryParam("dateFrom")
	dateTo := c.QueryParam("dateTo")
	if dateFrom != "" {
		fromDate, err := time.Parse("2006-01-02", dateFrom)
		if err == nil {
			query = query.Where("date >= ?", fromDate)
		}
	}
	if dateTo != "" {
		toDate, err := time.Parse("2006-01-02", dateTo)
		if err == nil {
			query = query.Where("date <= ?", toDate)
		}
	}

	var expenses []models.Expense
	if err := query.Order("date DESC, created_at DESC").Find(&expenses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to search expenses"})
	}

	// Format response with category names
	type ExpenseResponse struct {
		ID         uint      `json:"id"`
		Categories []string  `json:"categories"`
		Date       time.Time `json:"date"`
		Notes      string    `json:"notes"`
		Amount     float64   `json:"amount"`
		Category   string    `json:"category"` // First category for display
	}

	results := make([]ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		categoryName := ""
		if len(expense.Categories) > 0 {
			categoryName = expense.Categories[0]
		}
		results[i] = ExpenseResponse{
			ID:         expense.ID,
			Categories: expense.Categories,
			Date:       expense.Date,
			Notes:      expense.Notes,
			Amount:     expense.Amount,
			Category:   categoryName,
		}
	}

	return c.JSON(http.StatusOK, results)
}

func (h *ExpenseHandler) UpdateExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid expense ID"})
	}

	var req UpdateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var expense models.Expense
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&expense).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Expense not found"})
	}

	expense.Categories = req.Categories
	expense.Notes = req.Notes
	expense.Amount = req.Amount

	if err := h.db.Save(&expense).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update expense"})
	}

	return c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) DeleteExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid expense ID"})
	}

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Unscoped().Delete(&models.Expense{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete expense"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Expense deleted successfully"})
}
