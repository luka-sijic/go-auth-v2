package handler

import (
	"app/internal/models"
	"app/pkg/secret"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Cookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "based"
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Secure = false
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "write cookie")
}

func Me(c echo.Context) error {
	cookie, err := c.Cookie("access")
	if err != nil {
		log.Println("No access cookie:", err)
		return c.JSON(http.StatusInternalServerError, "cookie not set")
	}
	raw := cookie.Value
	log.Println("Raw JWT from cookie:", raw)

	token, err := jwt.ParseWithClaims(raw, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret.JwtSecret, nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token 1")
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		log.Println("Invalid tolen", token)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token or wrong type")
	}

	return c.JSON(http.StatusOK, claims.Username)
}

func Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		log.Println("No refresh cookie", err)
		return c.JSON(http.StatusInternalServerError, "no cookie")
	}
	raw := cookie.Value
	log.Println("Raw JWT from cookie", raw)

	token, err := jwt.ParseWithClaims(raw, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret.JwtSecret, nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token 1")
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		log.Println("Invalid tolen", token)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token or wrong type")
	}

	access := secret.GenerateJWT(claims.Username, 900)
	setCookie(c, "access", access, 900)

	return c.NoContent(http.StatusOK)
}

func setCookie(c echo.Context, name, val string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge * time.Now().UTC().Minute(),
	})
}
