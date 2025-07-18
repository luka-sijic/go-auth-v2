package routes

import (
	"app/internal/handler"
	"app/pkg/secret"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	e.GET("/me", handler.Me, secret.Auth)

	e.POST("/friend", handler.AddFriend, secret.Auth)
	e.GET("/friend", handler.GetRequest, secret.Auth)
	e.POST("/friend/respond", handler.Respond, secret.Auth)
	e.GET("/friend/:id", handler.GetFriends, secret.Auth)
}
