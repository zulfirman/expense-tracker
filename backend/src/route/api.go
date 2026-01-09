package route

import (
	"expenses-tracker/src/registry"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo, reg *registry.Registry) {
	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes (no authentication required)
	api := e.Group("/api/apps")
	api.POST("/auth/signup", reg.AuthHandler.Signup)
	api.POST("/auth/login", reg.AuthHandler.Login)
	api.POST("/auth/refresh", reg.AuthHandler.RefreshToken)

	// Protected routes (authentication required)
	protected := api.Group("")
	protected.Use(reg.AuthMiddleware)

	// Auth routes
	protected.GET("/auth/profile", reg.AuthHandler.GetProfile)
	protected.PUT("/auth/profile", reg.AuthHandler.UpdateProfile)
	protected.PUT("/auth/currency", reg.AuthHandler.UpdateCurrency)
	protected.PUT("/auth/password", reg.AuthHandler.ChangePassword)

	// Expense routes
	protected.POST("/expenses", reg.ExpenseHandler.CreateExpense)
	protected.GET("/expenses/months", reg.ExpenseHandler.GetMonths)
	protected.GET("/expenses/month/:month", reg.ExpenseHandler.GetMonthDetails)
	protected.GET("/expenses/date/:date", reg.ExpenseHandler.GetDateExpenses)
	protected.GET("/expenses/search", reg.ExpenseHandler.SearchExpenses)
	protected.PUT("/expenses/:id", reg.ExpenseHandler.UpdateExpense)
	protected.DELETE("/expenses/:id", reg.ExpenseHandler.DeleteExpense)

	// Category routes
	protected.GET("/categories", reg.CategoryHandler.GetCategories)
	protected.POST("/categories", reg.CategoryHandler.CreateCategory)
	protected.PUT("/categories/:id", reg.CategoryHandler.UpdateCategory)
	protected.PUT("/categories/sequence", reg.CategoryHandler.UpdateCategoriesSequence)
	protected.DELETE("/categories/:id", reg.CategoryHandler.DeleteCategory)

	// Template routes
	protected.GET("/templates", reg.TemplateHandler.GetTemplates)
	protected.POST("/templates", reg.TemplateHandler.CreateTemplate)
	protected.DELETE("/templates/:id", reg.TemplateHandler.DeleteTemplate)

	// Income routes
	protected.POST("/income", reg.IncomeHandler.CreateIncome)
	protected.GET("/income/balance", reg.IncomeHandler.GetBalance)
	protected.PUT("/income/balance", reg.IncomeHandler.UpdateBalance)
	protected.GET("/income/date/:date", reg.IncomeHandler.GetDateIncome)
	protected.PUT("/income/:id", reg.IncomeHandler.UpdateIncome)
	protected.DELETE("/income/:id", reg.IncomeHandler.DeleteIncome)

	// Budget routes
	protected.GET("/budgets", reg.BudgetHandler.GetBudgets)
	protected.POST("/budgets", reg.BudgetHandler.CreateBudget)
	protected.POST("/budgets/copy", reg.BudgetHandler.CopyBudgets)
	protected.GET("/budgets/latest", reg.BudgetHandler.GetLatestBudgetMonth)
	protected.DELETE("/budgets/:categoryId", reg.BudgetHandler.DeleteBudget)

	// Quick amounts routes
	protected.GET("/quick-amounts", reg.QuickAmountHandler.GetQuickAmounts)
	protected.PUT("/quick-amounts", reg.QuickAmountHandler.SetQuickAmounts)
}
