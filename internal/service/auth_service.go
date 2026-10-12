package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/mardenimihardi/employee-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidInput = errors.New("invalid registration input")

type AuthService struct {
	employees *repository.EmployeeRepository
}

func NewAuthService(
	employees *repository.EmployeeRepository,
) *AuthService {
	return &AuthService{employees: employees}
}

type RegisterInput struct {
	EmployeeNo     string
	DepartmentCode string
	PositionCode   string
	FullName       string
	Email          string
	Phone          string
	DateOfBirth    string
	Gender         string
	Address        string
	City           string
	Province       string
	PostalCode     string
	Password       string
}

func (s *AuthService) Register(
	ctx context.Context,
	input RegisterInput,
) error {
	input.EmployeeNo = strings.TrimSpace(input.EmployeeNo)
	input.DepartmentCode = strings.TrimSpace(input.DepartmentCode)
	input.PositionCode = strings.TrimSpace(input.PositionCode)
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Phone = strings.TrimSpace(input.Phone)
	input.Gender = strings.ToLower(strings.TrimSpace(input.Gender))
	input.Address = strings.TrimSpace(input.Address)
	input.City = strings.TrimSpace(input.City)
	input.Province = strings.TrimSpace(input.Province)
	input.PostalCode = strings.TrimSpace(input.PostalCode)

	if input.EmployeeNo == "" ||
		input.FullName == "" ||
		input.Email == "" ||
		input.Password == "" {
		return ErrInvalidInput
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return ErrInvalidInput
	}

	if input.Gender != "" &&
		input.Gender != "male" &&
		input.Gender != "female" {
		return ErrInvalidInput
	}

	if input.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil || dob.After(time.Now()) {
			return ErrInvalidInput
		}
	}

	if len(input.Password) < 8 || len([]byte(input.Password)) > 72 {
		return ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.employees.Create(ctx, repository.CreateEmployeeInput{
		EmployeeNo:     input.EmployeeNo,
		DepartmentCode: input.DepartmentCode,
		PositionCode:   input.PositionCode,
		FullName:       input.FullName,
		Email:          input.Email,
		Phone:          input.Phone,
		DateOfBirth:    input.DateOfBirth,
		Gender:         input.Gender,
		Address:        input.Address,
		City:           input.City,
		Province:       input.Province,
		PostalCode:     input.PostalCode,
		PasswordHash:   string(passwordHash),
	})
}
