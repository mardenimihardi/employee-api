package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mardenimihardi/employee-api/internal/service"
)

type LoginHandler struct {
	loginService *service.LoginService
}

func NewLoginHandler(loginService *service.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

type loginRequest struct {
	EmployeeNo string `json:"employee_no"`
	Password   string `json:"password"`
}

func (h *LoginHandler) Login(c echo.Context) error {
	var req loginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid JSON request",
		})
	}

	result, err := h.loginService.Login(
		c.Request().Context(),
		service.LoginInput{
			EmployeeNo: req.EmployeeNo,
			Password:   req.Password,
		},
	)

	switch {
	case err == nil:
		return c.JSON(http.StatusOK, result)

	case errors.Is(err, service.ErrInvalidCredentials):
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Employee number or password is incorrect",
		})

	case errors.Is(err, service.ErrAccountInactive):
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "Account is not active",
		})

	case errors.Is(err, service.ErrJWTConfiguration):
		c.Logger().Errorf("JWT configuration error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})

	case errors.Is(err, service.ErrAccountLocked):
		return c.JSON(http.StatusTooManyRequests, map[string]string{
			"message": "Too many failed login attempts. Try again later.",
		})

	default:
		c.Logger().Errorf("Login failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}
}
