package routes

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/secret"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo, svc *service.Infra) {
	userHandler := handler.NewHandler(svc)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	e.GET("/me", handler.Me)
	e.GET("/cookie", handler.Cookie)

	e.POST("/friend", userHandler.AddFriend, secret.Auth)
	e.GET("/friend", userHandler.GetRequest, secret.Auth)
	e.POST("/friend/respond", userHandler.Respond, secret.Auth)
	e.GET("/friend/:id", userHandler.GetFriends, secret.Auth)
	e.GET("/friend/:user1/:user2", userHandler.GetLog)
}
