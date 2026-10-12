package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mardenimihardi/employee-api/internal/config"
	"github.com/mardenimihardi/employee-api/internal/database"
	"github.com/mardenimihardi/employee-api/internal/handler"
	"github.com/mardenimihardi/employee-api/internal/repository"
	"github.com/mardenimihardi/employee-api/internal/routes"
	"github.com/mardenimihardi/employee-api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	employeeRepo := repository.NewEmployeeRepository(db)
	authService := service.NewAuthService(employeeRepo)
	authHandler := handler.NewAuthHandler(authService)
	meHandler := handler.NewMeHandler()
	loginService := service.NewLoginService(employeeRepo, cfg)
	loginHandler := handler.NewLoginHandler(loginService)
	attendanceService := service.NewAttendanceService(repository.NewAttendanceRepository(db))
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	routes.Register(e, authHandler, loginHandler, meHandler, attendanceHandler, cfg, db)

	log.Printf("Server running on port %s", cfg.AppPort)
	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
