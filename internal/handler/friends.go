package handler

import (
	"app/internal/models"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *UserHandler) AddFriend(c echo.Context) error {
	username := c.Get("username").(string)
	user := new(models.FriendDTO)
	if err := c.Bind(&user); err != nil {
		log.Println(err)
	}

	fmt.Println(username)
	fmt.Println(user.Friend)

	result := h.svc.AddFriend(username, user)
	if !result {
		return c.JSON(http.StatusNotFound, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend request sent")
}

func (h *UserHandler) GetFriends(c echo.Context) error {
	username := c.Param("id")

	result := h.svc.GetFriends(username)
	if len(result) == 0 {
		return c.JSON(http.StatusNotFound, "No friends found")
	}
	return c.JSON(http.StatusOK, result)
}

func (h *UserHandler) GetRequest(c echo.Context) error {
	username := c.Get("username").(string)

	result := h.svc.GetRequests(username)
	if len(result) == 0 {
		return c.JSON(http.StatusOK, "")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *UserHandler) Respond(c echo.Context) error {
	username := c.Get("username").(string)
	action := new(models.FriendActionDTO)
	if err := c.Bind(&action); err != nil {
		log.Println(err)
	}

	result := h.svc.FriendResponse(username, action)
	if !result {
		fmt.Println("Error")
		return c.JSON(http.StatusInternalServerError, "Failed to update friend request")
	}
	return c.JSON(http.StatusOK, "Friend request updated")
}
