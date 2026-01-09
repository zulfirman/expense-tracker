package handler

import (
	"net/http"
	"strconv"

	"expenses-tracker/src/middleware"
	"expenses-tracker/src/model"
	"expenses-tracker/src/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type WorkspaceHandler struct {
	repo     *repository.WorkspaceRepository
	userRepo *repository.UserRepository
	db       *gorm.DB
}

func NewWorkspaceHandler(repo *repository.WorkspaceRepository, userRepo *repository.UserRepository, db *gorm.DB) *WorkspaceHandler {
	return &WorkspaceHandler{
		repo:     repo,
		userRepo: userRepo,
		db:       db,
	}
}

type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (h *WorkspaceHandler) List(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	q := c.QueryParam("q")
	items, err := h.repo.ListByUser(userID, q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to list workspaces"})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *WorkspaceHandler) Get(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid workspace ID"})
	}

	ws, err := h.repo.GetByID(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Workspace not found"})
	}

	return c.JSON(http.StatusOK, ws)
}

func (h *WorkspaceHandler) Create(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var req CreateWorkspaceRequest
	if err := c.Bind(&req); err != nil || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	ws := model.M_workspace{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.Create(&ws); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create workspace"})
	}

	// Seed default categories for this workspace
	if err := repository.SeedDefaultCategoriesForWorkspace(h.db, userID, ws.ID); err != nil {
		// Log error but don't fail workspace creation
	}

	// Mark first signin as completed if this is the first workspace
	user, err := h.userRepo.GetByID(userID)
	if err == nil && !user.FirstSigninCompleted {
		user.FirstSigninCompleted = true
		h.userRepo.Update(user)
	}

	return c.JSON(http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Update(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid workspace ID"})
	}

	var req CreateWorkspaceRequest
	if err := c.Bind(&req); err != nil || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	ws, err := h.repo.GetByID(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Workspace not found"})
	}

	ws.Name = req.Name
	ws.Description = req.Description
	if err := h.repo.Update(ws); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update workspace"})
	}

	return c.JSON(http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid workspace ID"})
	}

	ws, err := h.repo.GetByID(userID, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Workspace not found"})
	}

	// Delete workspace (hard delete - cascade will delete related data)
	// Use Unscoped() to perform a hard delete instead of soft delete
	if err := h.db.Unscoped().Delete(ws).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete workspace: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace deleted successfully"})
}

