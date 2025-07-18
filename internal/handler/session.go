package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Me(c echo.Context) error {
	user := c.Get("username").(string)
	return c.JSON(http.StatusOK, map[string]string{"username": user})
}

func setCookie(c echo.Context, name, val string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		//SameSite: http.SameSiteLaxMode,
		MaxAge: maxAge,
	})
}
