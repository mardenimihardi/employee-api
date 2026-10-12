package routes

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mardenimihardi/employee-api/internal/config"
	"github.com/mardenimihardi/employee-api/internal/handler"
	appmiddleware "github.com/mardenimihardi/employee-api/internal/middleware"
)

func Register(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	loginHandler *handler.LoginHandler,
	meHandler *handler.MeHandler,
	attendanceHandler *handler.AttendanceHandler,
	cfg config.Config,
	db *sql.DB,
) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Employee API is running",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	api := e.Group("/api/v1")
	auth := api.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", loginHandler.Login)
	api.GET(
		"/auth/me",
		meHandler.Me,
		appmiddleware.JWTAuth(cfg, db),
	)
	api.POST(
		"/attendance",
		attendanceHandler.Record,
		appmiddleware.JWTAuth(cfg, db),
	)
}
