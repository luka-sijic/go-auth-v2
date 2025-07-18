package handler

import (
	"app/internal/models"
	"app/internal/service"
	"log"
	"net/http"

	"app/pkg/secret"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	user := new(models.UserDTO)
	if err := c.Bind(user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bind to user"})
	}

	result := service.CreateUser(user)
	if !result {
		return c.JSON(http.StatusInternalServerError, "Failed to register user")
	}

	return c.JSON(http.StatusOK, "User successfully registered")
}

func Login(c echo.Context) error {
	user := new(models.UserDTO)
	if err := c.Bind(&user); err != nil {
		log.Println(err)
	}

	result := service.LoginUser(user)
	if !result {
		return c.JSON(http.StatusInternalServerError, "Incorrect username/password")
	}

	access := secret.GenerateJWT(user.Username, 900)
	refresh := secret.GenerateJWT(user.Username, 2592000)

	setCookie(c, "access", access, 900)
	setCookie(c, "refresh", refresh, 2592000)

	return c.NoContent(http.StatusOK)
}
