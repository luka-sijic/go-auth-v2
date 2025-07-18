package handler

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"
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

func AddFriend(c echo.Context) error {
	username := c.Get("username").(string)
	user := new(models.FriendDTO)
	if err := c.Bind(&user); err != nil {
		log.Println(err)
	}

	fmt.Println(username)
	fmt.Println(user.Friend)

	result := service.AddFriend(username, user)
	if !result {
		return c.JSON(http.StatusNotFound, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend request sent")
}

func GetFriends(c echo.Context) error {
	username := c.Param("id")

	result := service.GetFriends(username)
	if len(result) == 0 {
		return c.JSON(http.StatusNotFound, "No friends found")
	}
	return c.JSON(http.StatusOK, result)
}

func GetRequest(c echo.Context) error {
	username := c.Get("username").(string)

	result := service.GetRequests(username)
	if len(result) == 0 {
		return c.JSON(http.StatusOK, "")
	}

	return c.JSON(http.StatusOK, result)
}

func Respond(c echo.Context) error {
	username := c.Get("username").(string)
	action := new(models.FriendActionDTO)
	if err := c.Bind(&action); err != nil {
		log.Println(err)
	}

	result := service.FriendResponse(username, action)
	if !result {
		fmt.Println("Error")
		return c.JSON(http.StatusInternalServerError, "Failed to update friend request")
	}
	return c.JSON(http.StatusOK, "Friend request updated")
}
