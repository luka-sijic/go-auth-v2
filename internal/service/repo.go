package service

import (
	"app/internal/database"
	"app/internal/models"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	CreateUser(user *models.UserDTO) bool
	LoginUser(user *models.UserDTO) bool

	AddFriend(username string, user *models.FriendDTO) bool
	GetFriends(username string) []models.FriendDTO
	GetRequests(username string) []models.FriendDTO
	FriendResponse(username string, action *models.FriendActionDTO) bool
	GetLog(username1, username2 string) []models.Messages
}

/*
type FrienService interface {
	AddFriend(username string, user *models.FriendDTO) bool
	GetFriends(username string) []models.FriendDTO
	GetRequests(username string) []models.FriendDTO
	FriendResponse(username string, action *models.FriendActionDTO) bool
}
*/

type Infra struct {
	Pools []*pgxpool.Pool
	Node  *snowflake.Node
	RDB   *redis.Client
}

func NewService(app *database.App) *Infra {
	return &Infra{Pools: app.Pools, Node: app.Node, RDB: app.RDB}
}
