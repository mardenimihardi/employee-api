package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrEmployeeExists = errors.New(
	"employee number or email already exists",
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

type EmployeeLogin struct {
	ID                  string
	EmployeeNo          string
	FullName            string
	Email               string
	PasswordHash        string
	AccountRole         string
	AccountStatus       string
	FailedLoginAttempts int
	LockedUntil         sql.NullTime
}

type CreateEmployeeInput struct {
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
	PasswordHash   string
}

func (r *EmployeeRepository) Create(
	ctx context.Context,
	input CreateEmployeeInput,
) error {
	query := `
        SELECT fn_register_employee(
            $1, $2, $3, $4,
            NULLIF($5, '')::DATE,
            $6, $7, $8, $9, $10,
            $11, $12, $13
        )
    `

	var employeeID string

	err := r.db.QueryRowContext(
		ctx,
		query,
		input.EmployeeNo,
		input.FullName,
		input.Email,
		input.Phone,
		input.DateOfBirth,
		input.Gender,
		input.Address,
		input.City,
		input.Province,
		input.PostalCode,
		input.PasswordHash,
		input.DepartmentCode,
		input.PositionCode,
	).Scan(&employeeID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmployeeExists
		}

		return err
	}

	return nil
}

func (r *EmployeeRepository) FindByEmployeeNo(
	ctx context.Context,
	employeeNo string,
) (*EmployeeLogin, error) {
	query := `SELECT * FROM fn_login_employee($1)`

	var employee EmployeeLogin

	err := r.db.QueryRowContext(ctx, query, employeeNo).Scan(
		&employee.ID,
		&employee.EmployeeNo,
		&employee.FullName,
		&employee.Email,
		&employee.PasswordHash,
		&employee.AccountRole,
		&employee.AccountStatus,
		&employee.FailedLoginAttempts,
		&employee.LockedUntil,
	)
	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepository) RecordFailedLogin(
	ctx context.Context,
	employeeID string,
) error {
	query := `
        UPDATE employees
        SET
            failed_login_attempts = CASE
                WHEN locked_until IS NOT NULL
                     AND locked_until <= NOW()
                THEN 1
                ELSE failed_login_attempts + 1
            END,
            locked_until = CASE
                WHEN (
                    CASE
                        WHEN locked_until IS NOT NULL
                             AND locked_until <= NOW()
                        THEN 1
                        ELSE failed_login_attempts + 1
                    END
                ) >= 5
                THEN NOW() + INTERVAL '15 minutes'
                ELSE NULL
            END
        WHERE id = $1
          AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, employeeID)
	return err
}

func (r *EmployeeRepository) ResetFailedLogin(
	ctx context.Context,
	employeeID string,
) error {
	query := `
        UPDATE employees
        SET
            failed_login_attempts = 0,
            locked_until = NULL
        WHERE id = $1
          AND deleted_at IS NULL
    `

	_, err := r.db.ExecContext(ctx, query, employeeID)
	return err
}

var _ = sql.ErrNoRows
