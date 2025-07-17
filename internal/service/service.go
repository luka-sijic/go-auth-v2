package service

import (
	"app/internal/database"
	"app/internal/models"
	"app/pkg/hash"
	"context"

	"log"

	"github.com/bwmarrin/snowflake"
)

func CreateUser(user *models.UserDTO) bool {
	res, err := database.RDB.Do(context.Background(), "BF.EXISTS", "users", user.Username).Bool()
	if err != nil || res {
		log.Println(err)
		return false
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		return false
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println(err)
		return false
	}

	id := node.Generate()
	db := database.GetShardPool(id)

	_, err = db.Exec(context.Background(), "INSERT INTO users (id, username, password) VALUES ($1,$2,$3)", id, user.Username, hashedPassword)
	if err != nil {
		log.Println("Failed to create user: ", err)
		return false
	}

	database.RDB.HSet(context.Background(), "user:id_map", user.Username, id.String())
	database.RDB.Do(context.Background(), "BF.ADD", "users", user.Username)

	return true
}

func LoginUser(user *models.UserDTO) bool {
	idStr, err := database.RDB.HGet(context.Background(), "user:id_map", user.Username).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		log.Println(err)
		return false
	}
	db := database.GetShardPool(id)

	var storedHash string
	err = db.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedHash)
	if err != nil || !hash.CheckPasswordHash(user.Password, storedHash) {
		log.Println(err)
		return false
	}

	log.Println(id)
	return true
}

func AddFriend(user *models.FriendDTO) bool {
	return false
}
