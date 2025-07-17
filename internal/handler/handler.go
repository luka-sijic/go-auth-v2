package handler

import (
	"app/internal/models"
	"app/internal/service"
	"log"
	"net/http"
	"time"

	"app/pkg/secret"

	"github.com/golang-jwt/jwt/v5"

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
	if err := c.Bind(user); err != nil {
		log.Println(err)
	}

	result := service.LoginUser(user)
	if !result {
		return c.JSON(http.StatusInternalServerError, "Failed to login user")
	}

	claims := &models.Claims{
		Username: user.Username,
		Role:     1,
		Status:   1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(secret.JwtSecret)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func AddFriend(c echo.Context) error {
	user := new(models.FriendDTO)
	if err := c.Bind(user); err != nil {
		log.Println(err)
	}

	result := service.AddFriend(user)
	if !result {
		return c.JSON(http.StatusInternalServerError, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend request sent")
}
