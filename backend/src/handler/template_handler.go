package handler

import (
	"net/http"
	"strconv"

	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
)

type TemplateHandler struct {
	templateRepo *repository.TemplateRepository
	categoryRepo *repository.CategoryRepository
}

func NewTemplateHandler(templateRepo *repository.TemplateRepository, categoryRepo *repository.CategoryRepository) *TemplateHandler {
	return &TemplateHandler{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *TemplateHandler) GetTemplates(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	items, err := h.templateRepo.GetByUser(cc.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch templates"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *TemplateHandler) CreateTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)

	var req struct {
		Name       string  `json:"name"`
		CategoryID uint    `json:"categoryId"`
		Amount     float64 `json:"amount"`
		Notes      string  `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	t := &model.M_expense_template{
		UserID: cc.UserID,
		Name:   req.Name,
		Amount: req.Amount,
		Notes:  req.Notes,
	}

	if err := h.templateRepo.Create(t); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create template"})
	}
	return c.JSON(http.StatusCreated, t)
}

func (h *TemplateHandler) DeleteTemplate(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	if err := h.templateRepo.Delete(uint(id), cc.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete template"})
	}
	return c.NoContent(http.StatusNoContent)
}
