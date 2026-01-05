package main

import (
	"expenses-tracker/database"
	"expenses-tracker/handlers"
	authMiddleware "expenses-tracker/middleware"
	"expenses-tracker/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.User{}, &models.Expense{}, &models.Category{}, &models.ExpenseTemplate{}, &models.Income{}, &models.Balance{}, &models.Budget{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Note: Categories are now user-specific, so we don't seed default categories
	// Each user will create their own categories

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	expenseHandler := handlers.NewExpenseHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	templateHandler := handlers.NewTemplateHandler(db)
	incomeHandler := handlers.NewIncomeHandler(db)
	budgetHandler := handlers.NewBudgetHandler(db)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes (no authentication required)
	api := e.Group("/api")
	api.POST("/auth/signup", authHandler.Signup)
	api.POST("/auth/login", authHandler.Login)

	// Protected routes (authentication required)
	protected := api.Group("")
	protected.Use(authMiddleware.CustomContextMiddleware(db))

	// Auth routes
	protected.GET("/auth/profile", authHandler.GetProfile)
	protected.PUT("/auth/profile", authHandler.UpdateProfile)

	// Expense routes
	protected.POST("/expenses", expenseHandler.CreateExpense)
	protected.GET("/expenses/months", expenseHandler.GetMonths)
	protected.GET("/expenses/month/:month", expenseHandler.GetMonthDetails)
	protected.GET("/expenses/date/:date", expenseHandler.GetDateExpenses)
	protected.GET("/expenses/search", expenseHandler.SearchExpenses)
	protected.PUT("/expenses/:id", expenseHandler.UpdateExpense)
	protected.DELETE("/expenses/:id", expenseHandler.DeleteExpense)

	// Category routes
	protected.GET("/categories", categoryHandler.GetCategories)
	protected.POST("/categories", categoryHandler.CreateCategory)
	protected.PUT("/categories/:id", categoryHandler.UpdateCategory)
	protected.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	// Template routes
	protected.GET("/templates", templateHandler.GetTemplates)
	protected.POST("/templates", templateHandler.CreateTemplate)
	protected.DELETE("/templates/:id", templateHandler.DeleteTemplate)

	// Income routes
	protected.POST("/income", incomeHandler.CreateIncome)
	protected.GET("/income/balance", incomeHandler.GetBalance)
	protected.PUT("/income/balance", incomeHandler.UpdateBalance)
	protected.GET("/income/date/:date", incomeHandler.GetDateIncome)
	protected.PUT("/income/:id", incomeHandler.UpdateIncome)
	protected.DELETE("/income/:id", incomeHandler.DeleteIncome)

	// Budget routes
	protected.GET("/budgets", budgetHandler.GetBudgets)
	protected.POST("/budgets", budgetHandler.CreateBudget)
	protected.DELETE("/budgets/:categorySlug", budgetHandler.DeleteBudget)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
