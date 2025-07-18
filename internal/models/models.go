package models

import "github.com/golang-jwt/jwt/v5"

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FriendDTO struct {
	Friend string `json:"friend"`
}

type FriendActionDTO struct {
	FriendID string `json:"friendId"`
	Action   string `json:"action"`
}

type Claims struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
	jwt.RegisteredClaims
}
