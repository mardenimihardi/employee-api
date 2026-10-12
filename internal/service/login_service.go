package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mardenimihardi/employee-api/internal/config"
	"github.com/mardenimihardi/employee-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid employee number or password")
	ErrAccountInactive    = errors.New("account is not active")
	ErrJWTConfiguration   = errors.New("JWT configuration is invalid")
	ErrAccountLocked      = errors.New("account is temporarily locked")
)

type LoginService struct {
	employees *repository.EmployeeRepository
	cfg       config.Config
}

func NewLoginService(
	employees *repository.EmployeeRepository,
	cfg config.Config,
) *LoginService {
	return &LoginService{
		employees: employees,
		cfg:       cfg,
	}
}

type LoginInput struct {
	EmployeeNo string
	Password   string
}

type LoginResult struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int64         `json:"expires_in"`
	Employee    LoginEmployee `json:"employee"`
}

type LoginEmployee struct {
	ID         string `json:"id"`
	EmployeeNo string `json:"employee_no"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

func (s *LoginService) Login(
	ctx context.Context,
	input LoginInput,
) (*LoginResult, error) {
	input.EmployeeNo = strings.TrimSpace(input.EmployeeNo)

	if input.EmployeeNo == "" ||
		input.Password == "" ||
		len([]byte(input.Password)) > 72 {
		return nil, ErrInvalidCredentials
	}

	employee, err := s.employees.FindByEmployeeNo(ctx, input.EmployeeNo)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if employee.LockedUntil.Valid &&
		employee.LockedUntil.Time.After(time.Now()) {
		return nil, ErrAccountLocked
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(employee.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		if recordErr := s.employees.RecordFailedLogin(
			ctx,
			employee.ID,
		); recordErr != nil {
			return nil, recordErr
		}

		return nil, ErrInvalidCredentials
	}

	if employee.AccountStatus != "active" {
		return nil, ErrAccountInactive
	}

	if err := s.employees.ResetFailedLogin(
		ctx,
		employee.ID,
	); err != nil {
		return nil, err
	}

	if strings.TrimSpace(s.cfg.JWTSecret) == "" {
		return nil, ErrJWTConfiguration
	}

	expiration := 15 * time.Minute
	if s.cfg.JWTExpiresIn != "" {
		parsed, err := time.ParseDuration(s.cfg.JWTExpiresIn)
		if err != nil || parsed <= 0 {
			return nil, ErrJWTConfiguration
		}
		expiration = parsed
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"sub":         employee.ID,
		"employee_no": employee.EmployeeNo,
		"role":        employee.AccountRole,
		"iat":         now.Unix(),
		"exp":         now.Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiration.Seconds()),
		Employee: LoginEmployee{
			ID:         employee.ID,
			EmployeeNo: employee.EmployeeNo,
			FullName:   employee.FullName,
			Email:      employee.Email,
			Role:       employee.AccountRole,
		},
	}, nil
}
