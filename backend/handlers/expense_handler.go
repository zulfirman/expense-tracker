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

// Helper function to enrich expense with category names
func (h *ExpenseHandler) enrichExpenseWithCategoryNames(expense models.T_expense) map[string]interface{} {
	categoryIDs := make([]uint, len(expense.Categories))
	categoryNames := make([]string, len(expense.Categories))
	for i, cat := range expense.Categories {
		categoryIDs[i] = cat.ID
		categoryNames[i] = cat.Name
	}

	return map[string]interface{}{
		"id":          expense.ID,
		"userId":      expense.UserID,
		"categoryIds": categoryIDs,
		"categories":  categoryNames,
		"date":        expense.Date,
		"notes":       expense.Notes,
		"amount":      expense.Amount,
		"createdAt":   expense.CreatedAt,
		"updatedAt":   expense.UpdatedAt,
	}
}

type CreateExpenseRequest struct {
	CategoryIDs []uint  `json:"categoryIds"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
	Amount      float64 `json:"amount"`
}

type UpdateExpenseRequest struct {
	CategoryIDs []uint  `json:"categoryIds"`
	Notes       string  `json:"notes"`
	Amount      float64 `json:"amount"`
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

	expense := models.T_expense{
		UserID:     userID,
		Categories: categories,
		Date:       date,
		Notes:      req.Notes,
		Amount:     req.Amount,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create expense"})
	}

	// Preload categories for response
	h.db.Preload("Categories").First(&expense, expense.ID)
	expenseResponse := h.enrichExpenseWithCategoryNames(expense)
	return c.JSON(http.StatusCreated, expenseResponse)
}

func (h *ExpenseHandler) GetMonths(c echo.Context) error {
	userID := middleware.GetCustomContext(c).UserID

	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Check if "before" parameter is provided for pagination
	beforeMonthStr := c.QueryParam("before")
	var startMonth time.Time
	var endMonth time.Time

	if beforeMonthStr != "" {
		// Parse the "before" month (YYYY-MM format)
		beforeMonth, err := time.Parse("2006-01", beforeMonthStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid before month format"})
		}
		// Load 6 months before the specified month
		endMonth = beforeMonth.AddDate(0, 0, -1) // End at the day before the beforeMonth
		startMonth = endMonth.AddDate(0, -5, 0)  // Start 6 months before
		// Adjust to first day of the month
		startMonth = time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, startMonth.Location())
		endMonth = time.Date(endMonth.Year(), endMonth.Month(), 1, 0, 0, 0, 0, endMonth.Location())
	} else {
		// Initial load: current month + 3 months back
		endMonth = currentMonth
		startMonth = currentMonth.AddDate(0, -3, 0)
	}

	type MonthResult struct {
		Month string  `json:"month"`
		Total float64 `json:"total"`
	}

	var results []MonthResult
	err := h.db.Model(&models.T_expense{}).
		Select("TO_CHAR(date, 'YYYY-MM') as month, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND date >= ? AND date < ?", userID, startMonth, endMonth.AddDate(0, 1, 0)).
		Group("TO_CHAR(date, 'YYYY-MM')").
		Order("month DESC").
		Scan(&results).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch months"})
	}

	// Generate all months in range (including empty ones)
	allMonths := make(map[string]bool)
	monthMap := make(map[string]float64)
	for _, result := range results {
		allMonths[result.Month] = true
		monthMap[result.Month] = result.Total
	}

	// Add all months in range (including empty ones)
	iterMonth := startMonth
	for !iterMonth.After(endMonth) {
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

	// Build response with all months (in descending order)
	var monthsWithDates []MonthWithDates
	iterMonth = endMonth
	for !iterMonth.Before(startMonth) {
		monthKey := iterMonth.Format("2006-01")

		type DateResult struct {
			Date  time.Time `json:"date"`
			Total float64   `json:"total"`
		}
		var dateResults []DateResult
		h.db.Model(&models.T_expense{}).
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
		h.db.Model(&models.T_income{}).
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
		if iterMonth.Before(startMonth) {
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
		CategoryID uint    `json:"categoryId"`
		Category   string  `json:"category"`
		Total      float64 `json:"total"`
	}

	type DailyTotal struct {
		Date     string  `json:"date"`
		Income   float64 `json:"income"`
		Expense  float64 `json:"expense"`
		NetTotal float64 `json:"netTotal"`
	}

	// Get all expenses for the month with categories preloaded
	var expenses []models.T_expense
	if err := h.db.Preload("Categories").Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).Find(&expenses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch month details"})
	}

	// Get all income for the month
	var incomes []models.T_income
	if err := h.db.Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).Find(&incomes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch income details"})
	}

	// Aggregate by category (keyed by category ID to keep a stable reference)
	type categoryAggregate struct {
		Name  string
		Total float64
	}
	categoryMap := make(map[uint]categoryAggregate)
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
			if agg, ok := categoryMap[category.ID]; ok {
				agg.Total += expense.Amount
				categoryMap[category.ID] = agg
			} else {
				categoryMap[category.ID] = categoryAggregate{
					Name:  category.Name,
					Total: expense.Amount,
				}
			}
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
	for id, agg := range categoryMap {
		categoryResults = append(categoryResults, CategoryTotal{
			CategoryID: id,
			Category:   agg.Name,
			Total:      agg.Total,
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

	var expenses []models.T_expense
	err = h.db.Preload("Categories").Where("user_id = ? AND date = ?", userID, date).Order("created_at DESC").Find(&expenses).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch expenses"})
	}

	// Enrich with category names
	enrichedExpenses := make([]map[string]interface{}, len(expenses))
	for i, expense := range expenses {
		enrichedExpenses[i] = h.enrichExpenseWithCategoryNames(expense)
	}

	return c.JSON(http.StatusOK, enrichedExpenses)
}

func (h *ExpenseHandler) SearchExpenses(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	query := h.db.Model(&models.T_expense{}).Where("user_id = ?", userID)

	// Search by text (notes)
	searchQuery := c.QueryParam("q")
	if searchQuery != "" {
		query = query.Where("notes ILIKE ?", "%"+searchQuery+"%")
	}

	// Filter by category ID using join table
	categoryIDStr := c.QueryParam("categoryId")
	if categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			query = query.Joins("JOIN t_expense_categories ON t_expenses.id = t_expense_categories.t_expense_id").
				Where("t_expense_categories.m_category_id = ? AND t_expenses.user_id = ?", categoryID, userID)
		}
	} else {
		// Ensure user_id filter is applied even without category filter
		query = query.Where("user_id = ?", userID)
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

	var expenses []models.T_expense
	if err := query.Preload("Categories").Order("date DESC, created_at DESC").Find(&expenses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to search expenses"})
	}

	// Enrich with category names
	results := make([]map[string]interface{}, len(expenses))
	for i, expense := range expenses {
		enriched := h.enrichExpenseWithCategoryNames(expense)
		// Add first category name for backward compatibility
		categories := enriched["categories"].([]string)
		if len(categories) > 0 {
			enriched["category"] = categories[0]
		} else {
			enriched["category"] = ""
		}
		results[i] = enriched
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

	var expense models.T_expense
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&expense).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Expense not found"})
	}

	expense.Categories = categories
	expense.Notes = req.Notes
	expense.Amount = req.Amount

	if err := h.db.Save(&expense).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update expense"})
	}

	// Preload categories for response
	h.db.Preload("Categories").First(&expense, expense.ID)
	expenseResponse := h.enrichExpenseWithCategoryNames(expense)
	return c.JSON(http.StatusOK, expenseResponse)
}

func (h *ExpenseHandler) DeleteExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid expense ID"})
	}

	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Unscoped().Delete(&models.T_expense{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete expense"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Expense deleted successfully"})
}
