package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	db           *gorm.DB
	expenseRepo  *repository.ExpenseRepository
	categoryRepo *repository.CategoryRepository
}

func NewExpenseHandler(db *gorm.DB, expenseRepo *repository.ExpenseRepository, categoryRepo *repository.CategoryRepository) *ExpenseHandler {
	return &ExpenseHandler{
		db:           db,
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *ExpenseHandler) CreateExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		CategoryIDs []uint  `json:"categoryIds"`
		Date        string  `json:"date"` // YYYY-MM-DD
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

	exp := &model.T_expense{
		UserID:     cc.UserID,
		Categories: cats,
		Date:       d,
		Notes:      req.Notes,
		Amount:     req.Amount,
	}

	if err := h.expenseRepo.Create(exp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create expense"})
	}
	return c.JSON(http.StatusCreated, exp)
}

type DateTotal struct {
	Date       string  `json:"date"`
	Total      float64 `json:"total"`
	HasIncome  bool    `json:"hasIncome,omitempty"`
	HasExpense bool    `json:"hasExpense,omitempty"`
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
	err := h.db.Model(&model.T_expense{}).
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
		Month    string      `json:"month"`
		Total    float64     `json:"total"`    // total expenses for the month
		NetTotal float64     `json:"netTotal"` // income - expenses for the month
		Dates    []DateTotal `json:"dates"`
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
		h.db.Model(&model.T_expense{}).
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
		h.db.Model(&model.T_income{}).
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

		// Convert to DateTotal format with income/expense flags and compute month net total
		monthNetTotal := 0.0
		dateTotals := make([]DateTotal, len(dateTotalsWithType))
		for j, dt := range dateTotalsWithType {
			monthNetTotal += dt.Total
			dateTotals[j] = DateTotal{
				Date:       dt.Date,
				Total:      dt.Total,
				HasIncome:  dt.HasIncome,
				HasExpense: dt.HasExpense,
			}
		}

		// Use expense total only for Total (used by remaining-balance calc)
		expenseTotal := monthMap[monthKey]

		monthsWithDates = append(monthsWithDates, MonthWithDates{
			Month:    monthKey,
			Total:    expenseTotal,
			NetTotal: monthNetTotal,
			Dates:    dateTotals,
		})

		iterMonth = iterMonth.AddDate(0, -1, 0)
		if iterMonth.Before(startMonth) {
			break
		}
	}

	return c.JSON(http.StatusOK, monthsWithDates)
}

func (h *ExpenseHandler) GetMonthDetails(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	month := c.Param("month") // YYYY-MM
	userID := cc.UserID

	// Load all expenses for this month with categories
	var expenses []model.T_expense
	if err := h.db.
		Preload("Categories").
		Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).
		Find(&expenses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch expenses"})
	}

	// Load all incomes for this month
	var incomes []model.T_income
	if err := h.db.
		Where("user_id = ? AND TO_CHAR(date, 'YYYY-MM') = ?", userID, month).
		Find(&incomes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch income"})
	}

	// Aggregate by category
	type CategorySummary struct {
		CategoryID uint    `json:"categoryId"`
		Category   string  `json:"category"`
		Total      float64 `json:"total"`
	}
	categoryMap := make(map[uint]*CategorySummary)

	for _, e := range expenses {
		for _, cat := range e.Categories {
			cs, ok := categoryMap[cat.ID]
			if !ok {
				cs = &CategorySummary{
					CategoryID: cat.ID,
					Category:   cat.Name,
					Total:      0,
				}
				categoryMap[cat.ID] = cs
			}
			cs.Total += e.Amount
		}
	}

	categories := make([]CategorySummary, 0, len(categoryMap))
	for _, cs := range categoryMap {
		categories = append(categories, *cs)
	}

	// Aggregate by day for income/expenses
	type DailySummary struct {
		Date     string  `json:"date"`
		Income   float64 `json:"income"`
		Expense  float64 `json:"expense"`
		NetTotal float64 `json:"netTotal"`
	}
	dailyMap := make(map[string]*DailySummary)

	for _, e := range expenses {
		key := e.Date.Format("2006-01-02")
		ds, ok := dailyMap[key]
		if !ok {
			ds = &DailySummary{Date: key}
			dailyMap[key] = ds
		}
		ds.Expense += e.Amount
	}
	for _, in := range incomes {
		key := in.Date.Format("2006-01-02")
		ds, ok := dailyMap[key]
		if !ok {
			ds = &DailySummary{Date: key}
			dailyMap[key] = ds
		}
		ds.Income += in.Amount
	}

	daily := make([]DailySummary, 0, len(dailyMap))
	for _, ds := range dailyMap {
		ds.NetTotal = ds.Income - ds.Expense
		daily = append(daily, *ds)
	}

	// Sort daily ascending by date
	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Date < daily[j].Date
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categories,
		"daily":      daily,
	})
}

func (h *ExpenseHandler) GetDateExpenses(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	dateStr := c.Param("date")
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date"})
	}

	items, err := h.expenseRepo.GetByDate(cc.UserID, d)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch expenses"})
	}

	// Shape for DateExpensesModal: categories as []string and categoryIds for editing
	type respItem struct {
		ID          uint      `json:"id"`
		Categories  []string  `json:"categories"`
		CategoryIDs []uint    `json:"categoryIds"`
		Date        time.Time `json:"date"`
		Notes       string    `json:"notes"`
		Amount      float64   `json:"amount"`
	}
	out := make([]respItem, 0, len(items))
	for _, e := range items {
		names := make([]string, 0, len(e.Categories))
		ids := make([]uint, 0, len(e.Categories))
		for _, c := range e.Categories {
			names = append(names, c.Name)
			ids = append(ids, c.ID)
		}
		out = append(out, respItem{
			ID:          e.ID,
			Categories:  names,
			CategoryIDs: ids,
			Date:        e.Date,
			Notes:       e.Notes,
			Amount:      e.Amount,
		})
	}

	return c.JSON(http.StatusOK, out)
}

func (h *ExpenseHandler) SearchExpenses(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	q := c.QueryParam("q")

	// Optional filters from query params
	var categoryID *uint
	if cidStr := c.QueryParam("categoryId"); cidStr != "" {
		if cid, err := strconv.Atoi(cidStr); err == nil && cid > 0 {
			val := uint(cid)
			categoryID = &val
		}
	}

	var dateFrom, dateTo *time.Time
	if v := c.QueryParam("dateFrom"); v != "" {
		if d, err := time.Parse("2006-01-02", v); err == nil {
			dateFrom = &d
		}
	}
	if v := c.QueryParam("dateTo"); v != "" {
		if d, err := time.Parse("2006-01-02", v); err == nil {
			dateTo = &d
		}
	}

	items, err := h.expenseRepo.Search(cc.UserID, q, categoryID, dateFrom, dateTo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to search expenses"})
	}

	// Shape response to what the frontend expects in ExpensesHistory.svelte
	type respItem struct {
		Date     time.Time `json:"date"`
		Category string    `json:"category"`
		Notes    string    `json:"notes"`
		Amount   float64   `json:"amount"`
	}
	out := make([]respItem, 0, len(items))
	for _, e := range items {
		// For multi-category expenses, join category names with comma
		categoryName := ""
		if len(e.Categories) > 0 {
			names := make([]string, 0, len(e.Categories))
			for _, c := range e.Categories {
				names = append(names, c.Name)
			}
			categoryName = strings.Join(names, ", ")
		}
		out = append(out, respItem{
			Date:     e.Date,
			Category: categoryName,
			Notes:    e.Notes,
			Amount:   e.Amount,
		})
	}

	return c.JSON(http.StatusOK, out)
}

func (h *ExpenseHandler) UpdateExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	exp, err := h.expenseRepo.GetByID(uint(id), cc.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Expense not found"})
	}

	var req struct {
		CategoryIDs *[]uint  `json:"categoryIds"`
		Date        *string  `json:"date"`
		Notes       *string  `json:"notes"`
		Amount      *float64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if req.CategoryIDs != nil {
		ids := *req.CategoryIDs
		cats := make([]model.M_category, len(ids))
		for i, id := range ids {
			cats[i] = model.M_category{ID: id}
		}
		// Replace categories via repository helper
		if err := h.expenseRepo.ReplaceCategories(exp, cats); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update categories"})
		}
	}
	if req.Date != nil {
		if d, err := time.Parse("2006-01-02", *req.Date); err == nil {
			exp.Date = d
		}
	}
	if req.Notes != nil {
		exp.Notes = *req.Notes
	}
	if req.Amount != nil {
		exp.Amount = *req.Amount
	}

	if err := h.expenseRepo.Update(exp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update expense"})
	}
	return c.JSON(http.StatusOK, exp)
}

func (h *ExpenseHandler) DeleteExpense(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	if err := h.expenseRepo.Delete(uint(id), cc.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete expense"})
	}
	return c.NoContent(http.StatusNoContent)
}
