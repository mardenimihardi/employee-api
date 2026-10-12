package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mardenimihardi/employee-api/internal/service"
)

type AttendanceHandler struct {
	attendanceService *service.AttendanceService
}

func NewAttendanceHandler(
	attendanceService *service.AttendanceService,
) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendanceService,
	}
}

type attendanceRequest struct {
	AttendanceType     string    `json:"attendance_type"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	AccuracyM          *float64  `json:"accuracy_m"`
	LocationCapturedAt time.Time `json:"location_captured_at"`
}

func (h *AttendanceHandler) Record(c echo.Context) error {
	var req attendanceRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	// Employee ID diambil dari JWT middleware,
	// bukan dari request body.
	employeeID, ok := c.Get("employee_id").(string)
	if !ok || employeeID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}

	result, err := h.attendanceService.Record(
		c.Request().Context(),
		service.AttendanceInput{
			EmployeeID:         employeeID,
			AttendanceType:     req.AttendanceType,
			Latitude:           req.Latitude,
			Longitude:          req.Longitude,
			AccuracyM:          req.AccuracyM,
			LocationCapturedAt: req.LocationCapturedAt,
		},
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLocationTimestamp) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid location timestamp",
				"detail":  "location_captured_at must be within the last 10 minutes and no more than 1 minute in the future",
			})
		}

		if errors.Is(err, service.ErrInvalidAttendance) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid attendance data",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to record attendance",
			"detail":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Attendance recorded successfully",
		"data":    result,
	})
}
