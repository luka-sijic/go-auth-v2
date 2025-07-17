package routes

import (
	"app/internal/handler"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
	e.POST("/addfriend", handler.AddFriend)
}
