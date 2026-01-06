package registry

import (
	"expenses-tracker/src/config"
	"expenses-tracker/src/handler"
	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Registry struct {
	DB *gorm.DB

	// Repositories
	UserRepo         *repository.UserRepository
	ExpenseRepo      *repository.ExpenseRepository
	IncomeRepo       *repository.IncomeRepository
	CategoryRepo     *repository.CategoryRepository
	BudgetRepo       *repository.BudgetRepository
	TemplateRepo     *repository.TemplateRepository
	RefreshTokenRepo *repository.RefreshTokenRepository
	QuickAmountRepo  *repository.QuickAmountRepository

	// Handlers
	AuthHandler        *handler.AuthHandler
	ExpenseHandler     *handler.ExpenseHandler
	IncomeHandler      *handler.IncomeHandler
	CategoryHandler    *handler.CategoryHandler
	BudgetHandler      *handler.BudgetHandler
	TemplateHandler    *handler.TemplateHandler
	QuickAmountHandler *handler.QuickAmountHandler

	// Middleware
	AuthMiddleware echo.MiddlewareFunc
}

func NewRegistry() (*Registry, error) {
	// Connect to database
	db, err := config.ConnectPostgreSQL()
	if err != nil {
		return nil, err
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&model.M_user{},
		&model.T_expense{},
		&model.M_category{},
		&model.M_expense_template{},
		&model.T_income{},
		&model.R_balance{},
		&model.R_budget{},
		&model.M_refresh_token{},
		&model.M_quick_amount{},
	); err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	quickAmountRepo := repository.NewQuickAmountRepository(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userRepo, refreshTokenRepo)
	expenseHandler := handler.NewExpenseHandler(db, expenseRepo, categoryRepo)
	incomeHandler := handler.NewIncomeHandler(incomeRepo)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	budgetHandler := handler.NewBudgetHandler(budgetRepo, categoryRepo)
	templateHandler := handler.NewTemplateHandler(templateRepo, categoryRepo)
	quickAmountHandler := handler.NewQuickAmountHandler(quickAmountRepo)

	// Initialize middleware
	authMiddleware := middleware.CustomContextMiddleware(userRepo)

	return &Registry{
		DB:                 db,
		UserRepo:           userRepo,
		ExpenseRepo:        expenseRepo,
		IncomeRepo:         incomeRepo,
		CategoryRepo:       categoryRepo,
		BudgetRepo:         budgetRepo,
		TemplateRepo:       templateRepo,
		RefreshTokenRepo:   refreshTokenRepo,
		QuickAmountRepo:    quickAmountRepo,
		AuthHandler:        authHandler,
		ExpenseHandler:     expenseHandler,
		IncomeHandler:      incomeHandler,
		CategoryHandler:    categoryHandler,
		BudgetHandler:      budgetHandler,
		TemplateHandler:    templateHandler,
		QuickAmountHandler: quickAmountHandler,
		AuthMiddleware:     authMiddleware,
	}, nil
}
