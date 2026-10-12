package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/mardenimihardi/employee-api/internal/repository"
)

var (
	ErrInvalidAttendance        = errors.New("invalid attendance data")
	ErrInvalidLocationTimestamp = errors.New(
		"location_captured_at must be within the last 10 minutes and no more than 1 minute in the future",
	)
)

type AttendanceService struct {
	attendanceRepo *repository.AttendanceRepository
}

func NewAttendanceService(
	attendanceRepo *repository.AttendanceRepository,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
	}
}

type AttendanceInput struct {
	EmployeeID         string
	AttendanceType     string
	Latitude           float64
	Longitude          float64
	AccuracyM          *float64
	LocationCapturedAt time.Time
}

func (s *AttendanceService) Record(
	ctx context.Context,
	input AttendanceInput,
) (*repository.AttendanceRecord, error) {
	input.EmployeeID = strings.TrimSpace(input.EmployeeID)
	input.AttendanceType = strings.ToLower(
		strings.TrimSpace(input.AttendanceType),
	)

	if input.EmployeeID == "" {
		return nil, ErrInvalidAttendance
	}

	if input.AttendanceType != "check_in" &&
		input.AttendanceType != "check_out" {
		return nil, ErrInvalidAttendance
	}

	if math.IsNaN(input.Latitude) ||
		math.IsInf(input.Latitude, 0) ||
		input.Latitude < -90 ||
		input.Latitude > 90 {
		return nil, ErrInvalidAttendance
	}

	if math.IsNaN(input.Longitude) ||
		math.IsInf(input.Longitude, 0) ||
		input.Longitude < -180 ||
		input.Longitude > 180 {
		return nil, ErrInvalidAttendance
	}

	if input.AccuracyM != nil &&
		(math.IsNaN(*input.AccuracyM) ||
			math.IsInf(*input.AccuracyM, 0) ||
			*input.AccuracyM < 0) {
		return nil, ErrInvalidAttendance
	}

	if input.LocationCapturedAt.IsZero() {
		return nil, ErrInvalidLocationTimestamp
	}

	now := time.Now()

	if input.LocationCapturedAt.After(now.Add(1*time.Minute)) ||
		input.LocationCapturedAt.Before(now.Add(-10*time.Minute)) {
		return nil, ErrInvalidLocationTimestamp
	}

	return s.attendanceRepo.Create(ctx, repository.CreateAttendanceInput{
		EmployeeID:         input.EmployeeID,
		AttendanceType:     input.AttendanceType,
		Latitude:           input.Latitude,
		Longitude:          input.Longitude,
		AccuracyM:          input.AccuracyM,
		LocationCapturedAt: input.LocationCapturedAt,
	})
}
