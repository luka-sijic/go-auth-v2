package server

import (
	"app/internal/database"
	"app/internal/routes"
	"app/internal/service"
	"app/pkg/auto"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app, err := database.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := auto.Init(app.Pools); err != nil {
		log.Fatalf("could not create users table: %v", err)
	}

	userSvc := service.NewService(app)
	routes.Routes(e, userSvc)

	e.Logger.Fatal(e.Start(":8081"))
}
