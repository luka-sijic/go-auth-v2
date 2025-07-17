package secret

import (
	"app/internal/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtToken string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	jwtToken = os.Getenv("JWT_TOKEN")
}

var JwtSecret = []byte(jwtToken)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "User must login")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			if claims.Status == -1 || claims.Status == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, "User is banned")
			}
			// Store the claims in the context for later use
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Set("status", claims.Status)
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
}

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "User must login")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		// Extract claims
		claims, ok := token.Claims.(*models.Claims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
		}
		//log.Println(claims.Role)
		log.Println(claims.Username)
		if claims.Role != 3 {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("status", claims.Status)

		return next(c)
	}
}
