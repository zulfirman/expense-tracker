package handler

import (
	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"
	"expenses-tracker/src/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryHandler(categoryRepo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=income expense"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type" validate:"omitempty,oneof=income expense"`
	IsActive bool   `json:"isActive"`
	Sequence int    `json:"sequence"`
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	// Optional type filter query parameter
	typeFilter := c.QueryParam("type")

	categories, err := h.categoryRepo.GetAll(userID, typeFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch categories"})
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Generate slug from name
	baseSlug := utils.GenerateSlug(req.Name)

	// Check if slug exists and generate unique one if needed
	slug := baseSlug
	slugExists, err := h.categoryRepo.SlugExists(userID, slug, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to validate slug"})
	}

	if slugExists {
		// Generate unique slug
		counter := 1
		for slugExists {
			slug = baseSlug + "-" + strconv.Itoa(counter)
			slugExists, err = h.categoryRepo.SlugExists(userID, slug, 0)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to validate slug"})
			}
			counter++
		}
	}

	category := model.M_category{
		UserID:   userID,
		Name:     req.Name,
		Slug:     slug,
		Type:     req.Type,
		IsActive: true,
	}

	if err := h.categoryRepo.Create(&category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create category"})
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	categoryID, err := parseUint(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	category, err := h.categoryRepo.GetByID(userID, categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Category not found"})
	}

	if req.Name != "" && req.Name != category.Name {
		// Generate new slug if name changed
		baseSlug := utils.GenerateSlug(req.Name)

		// Check if slug exists (excluding current category)
		slug := baseSlug
		slugExists, err := h.categoryRepo.SlugExists(userID, slug, categoryID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to validate slug"})
		}

		if slugExists {
			// Generate unique slug
			counter := 1
			for slugExists {
				slug = baseSlug + "-" + strconv.Itoa(counter)
				slugExists, err = h.categoryRepo.SlugExists(userID, slug, categoryID)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to validate slug"})
				}
				counter++
			}
		}

		category.Name = req.Name
		category.Slug = slug
	}
	if req.Type != "" && (req.Type == "income" || req.Type == "expense") {
		category.Type = req.Type
	}
	if req.IsActive != category.IsActive {
		category.IsActive = req.IsActive
	}
	// Update sequence if provided (0 is valid)
	if req.Sequence != 0 || (req.Sequence == 0 && category.Sequence != 0) {
		category.Sequence = req.Sequence
	}

	if err := h.categoryRepo.Update(category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update category"})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID
	id := c.Param("id")

	categoryID, err := parseUint(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid category ID"})
	}

	if err := h.categoryRepo.Delete(userID, categoryID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Category not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) UpdateCategoriesSequence(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req struct {
		Categories []struct {
			ID       uint `json:"id"`
			Sequence int  `json:"sequence"`
		} `json:"categories"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Update each category's sequence
	for _, catReq := range req.Categories {
		category, err := h.categoryRepo.GetByID(userID, catReq.ID)
		if err != nil {
			continue // Skip if category not found
		}
		category.Sequence = catReq.Sequence
		if err := h.categoryRepo.Update(category); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update category sequence"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Categories reordered successfully"})
}

func parseUint(s string) (uint, error) {
	result, err := strconv.ParseUint(s, 10, 32)
	return uint(result), err
}
