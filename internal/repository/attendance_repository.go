package repository

import (
	"context"
	"database/sql"
	"time"
)

type AttendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

type CreateAttendanceInput struct {
	EmployeeID         string
	AttendanceType     string
	Latitude           float64
	Longitude          float64
	AccuracyM          *float64
	LocationCapturedAt time.Time
}

type AttendanceRecord struct {
	ID                 string    `json:"id"`
	EmployeeID         string    `json:"employee_id"`
	AttendanceType     string    `json:"attendance_type"`
	RecordedAt         time.Time `json:"recorded_at"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	AccuracyM          float64   `json:"accuracy_m"`
	LocationCapturedAt time.Time `json:"location_captured_at"`
}

func (r *AttendanceRepository) Create(
	ctx context.Context,
	input CreateAttendanceInput,
) (*AttendanceRecord, error) {
	query := `
		SELECT
			id,
			employee_id,
			attendance_type,
			recorded_at,
			latitude,
			longitude,
			accuracy_m,
			location_captured_at
		FROM fn_record_attendance($1, $2, $3, $4, $5, $6)
	`

	var record AttendanceRecord

	err := r.db.QueryRowContext(
		ctx,
		query,
		input.EmployeeID,
		input.AttendanceType,
		input.Latitude,
		input.Longitude,
		input.AccuracyM,
		input.LocationCapturedAt,
	).Scan(
		&record.ID,
		&record.EmployeeID,
		&record.AttendanceType,
		&record.RecordedAt,
		&record.Latitude,
		&record.Longitude,
		&record.AccuracyM,
		&record.LocationCapturedAt,
	)
	if err != nil {
		return nil, err
	}

	return &record, nil
}
