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
		AllowOrigins: []string{"http://127.0.0.1:3000", "http://localhost:3000"},
		AllowMethods: []string{echo.OPTIONS, echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderCookie,
		},
		ExposeHeaders:    []string{echo.HeaderSetCookie},
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
