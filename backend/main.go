package main

import (
	"errors"
	"expenses-tracker/database"
	"expenses-tracker/handlers"
	"expenses-tracker/models"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
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
	if err := db.AutoMigrate(&models.Expense{}, &models.Category{}, &models.ExpenseTemplate{}, &models.Income{}, &models.Balance{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default categories if they don't exist
	seedCategories(db)

	// Initialize handlers
	expenseHandler := handlers.NewExpenseHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	templateHandler := handlers.NewTemplateHandler(db)
	incomeHandler := handlers.NewIncomeHandler(db)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API routes
	api := e.Group("/api")
	api.POST("/expenses", expenseHandler.CreateExpense)
	api.GET("/expenses/months", expenseHandler.GetMonths)
	api.GET("/expenses/month/:month", expenseHandler.GetMonthDetails)
	api.GET("/expenses/date/:date", expenseHandler.GetDateExpenses)
	api.PUT("/expenses/:id", expenseHandler.UpdateExpense)
	api.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
	api.GET("/categories", categoryHandler.GetCategories)
	api.GET("/templates", templateHandler.GetTemplates)
	api.POST("/templates", templateHandler.CreateTemplate)
	api.DELETE("/templates/:id", templateHandler.DeleteTemplate)
	api.POST("/income", incomeHandler.CreateIncome)
	api.GET("/income/balance", incomeHandler.GetBalance)
	api.PUT("/income/balance", incomeHandler.UpdateBalance)
	api.GET("/income/date/:date", incomeHandler.GetDateIncome)
	api.PUT("/income/:id", incomeHandler.UpdateIncome)
	api.DELETE("/income/:id", incomeHandler.DeleteIncome)

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

func seedCategories(db *gorm.DB) {
	defaultCategories := []models.Category{
		{Name: "Makanan & Minuman", Slug: "makanan-minuman", IsActive: true},
		{Name: "Transportasi", Slug: "transportasi", IsActive: true},
		{Name: "Bensin", Slug: "bensin", IsActive: true},
		{Name: "Tagihan Listrik", Slug: "tagihan-listrik", IsActive: true},
		{Name: "Kontrakan", Slug: "kontrakan", IsActive: true},
		{Name: "Buat Ayang", Slug: "buat-ayang", IsActive: true},
		{Name: "Internet & Pulsa", Slug: "internet-pulsa", IsActive: true},
		{Name: "Belanja Bulanan", Slug: "belanja-bulanan", IsActive: true},
		{Name: "Jalan-jalan", Slug: "jalan-jalan", IsActive: true},
		{Name: "Tabungan", Slug: "tabungan", IsActive: true},
		{Name: "Lain-lain", Slug: "lain-lain", IsActive: true},
	}

	for _, category := range defaultCategories {
		var existing models.Category
		if err := db.Where("slug = ?", category.Slug).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Fallback: auto-generate slug if empty (though we already set them)
				if category.Slug == "" {
					category.Slug = strings.ToLower(strings.ReplaceAll(category.Name, " ", "-"))
				}

				if err := db.Create(&category).Error; err != nil {
					log.Printf("Failed to create category %s: %v", category.Name, err)
				} else {
					log.Printf("Created default category: %s", category.Name)
				}
			}
		} else {
			// Optional: ensure it's active even if exists
			if !existing.IsActive {
				db.Model(&existing).Update("is_active", true)
			}
		}
	}
}
