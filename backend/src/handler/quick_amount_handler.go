package handler

import (
	"expenses-tracker/src/middleware"
	"expenses-tracker/src/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type QuickAmountHandler struct {
	repo *repository.QuickAmountRepository
}

func NewQuickAmountHandler(repo *repository.QuickAmountRepository) *QuickAmountHandler {
	return &QuickAmountHandler{repo: repo}
}

type QuickAmountPayload struct {
	Amounts []float64 `json:"amounts"`
}

func (h *QuickAmountHandler) GetQuickAmounts(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	list, err := h.repo.GetByUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to load quick amounts"})
	}

	out := make([]float64, len(list))
	for i, qa := range list {
		out[i] = qa.Value
	}

	return c.JSON(http.StatusOK, out)
}

func (h *QuickAmountHandler) SetQuickAmounts(c echo.Context) error {
	cc := middleware.GetCustomContext(c)
	userID := cc.UserID

	var payload QuickAmountPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid payload"})
	}
	amounts := make([]float64, 0, len(payload.Amounts))
	for _, v := range payload.Amounts {
		if v > 0 {
			amounts = append(amounts, v)
		}
	}
	if err := h.repo.ReplaceForUser(userID, amounts); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save quick amounts"})
	}
	return c.NoContent(http.StatusOK)
}
