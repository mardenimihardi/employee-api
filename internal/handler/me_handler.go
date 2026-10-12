package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MeHandler struct{}

func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

func (h *MeHandler) Me(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Authenticated successfully",
		"employee_id": c.Get("employee_id"),
	})
}
