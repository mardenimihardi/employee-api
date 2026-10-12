package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mardenimihardi/employee-api/internal/repository"
	"github.com/mardenimihardi/employee-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	EmployeeNo     string `json:"employee_no"`
	DepartmentCode string `json:"department_code"`
	PositionCode   string `json:"position_code"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Province       string `json:"province"`
	PostalCode     string `json:"postal_code"`
	Password       string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req registerRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid JSON request",
		})
	}

	err := h.authService.Register(
		c.Request().Context(),
		service.RegisterInput{
			EmployeeNo:     req.EmployeeNo,
			DepartmentCode: req.DepartmentCode,
			PositionCode:   req.PositionCode,
			FullName:       req.FullName,
			Email:          req.Email,
			Phone:          req.Phone,
			DateOfBirth:    req.DateOfBirth,
			Gender:         req.Gender,
			Address:        req.Address,
			City:           req.City,
			Province:       req.Province,
			PostalCode:     req.PostalCode,
			Password:       req.Password,
		},
	)

	switch {
	case err == nil:
		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Employee registered successfully",
		})

	case errors.Is(err, service.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid registration input. Check required fields, email, gender, date of birth, and password.",
		})

	case errors.Is(err, repository.ErrEmployeeExists):
		return c.JSON(http.StatusConflict, map[string]string{
			"message": "Employee number or email already exists",
		})

	default:
		c.Logger().Errorf("Registration failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}
}
