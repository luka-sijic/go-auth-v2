package service

import (
	"app/internal/database"
	"app/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

func (infra *Infra) AddFriend(username string, user *models.FriendDTO) bool {
	res, err := infra.RDB.Do(context.Background(), "BF.EXISTS", "users", user.Friend).Bool()
	if err != nil || !res {
		log.Printf("User not found: %v", err)
		return false
	}

	val, err := infra.RDB.HMGet(context.Background(), "user:id_map", username, user.Friend).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	sender := val[0].(string)
	receiver := val[1].(string)
	fmt.Printf("TESTING: %s %s", sender, receiver)

	id, err := snowflake.ParseString(sender)
	if err != nil {
		log.Println(err)
		return false
	}
	db := database.GetShardPool(infra.Pools, id)

	_, err = db.Exec(context.Background(), "INSERT INTO friends (user_id, status, requester_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING", receiver, "pending", id)
	if err != nil {
		log.Println("Failed to create user: ", err)
		return false
	}

	id, err = snowflake.ParseString(receiver)
	if err != nil {
		log.Println(err)
		return false
	}
	db = database.GetShardPool(infra.Pools, id)

	_, err = db.Exec(context.Background(), "INSERT INTO friends (user_id, status, requester_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING", receiver, "pending", sender)
	if err != nil {
		log.Println("Failed to create user: ", err)
		return false
	}

	return true
}

func (infra *Infra) GetFriends(username string) []models.FriendDTO {
	val, err := infra.RDB.HGet(context.Background(), "user:id_map", username).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	id, err := snowflake.ParseString(val)
	if err != nil {
		log.Println(err)
		return nil
	}
	db := database.GetShardPool(infra.Pools, id)
	var friends []models.FriendDTO
	const sql = `
		SELECT user_id
		FROM friends 
		WHERE requester_id=$1 
		AND status='accepted'
		UNION
		SELECT requester_id
		FROM friends
		WHERE user_id=$1
		AND status='accepted'
	`
	rows, err := db.Query(context.Background(), sql, id.Int64())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var friend int64
		err := rows.Scan(&friend)
		if err != nil {
			log.Println(err)
			return nil
		}
		friendid := fmt.Sprintf("%d", friend)
		val, err := infra.RDB.HGet(context.Background(), "user:username_map", friendid).Result()
		if err != nil {
			log.Println(err)
			return nil
		}
		friendUsername := models.FriendDTO{Friend: val}
		friends = append(friends, friendUsername)
	}
	return friends
}

func (infra *Infra) GetRequests(username string) []models.FriendDTO {
	val, err := infra.RDB.HGet(context.Background(), "user:id_map", username).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	id, err := snowflake.ParseString(val)
	if err != nil {
		log.Println(err)
		return nil
	}
	db := database.GetShardPool(infra.Pools, id)
	var friends []models.FriendDTO
	rows, err := db.Query(context.Background(), "SELECT requester_id FROM friends WHERE user_id=$1 AND status='pending'", id.Int64())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var friend int64
		err := rows.Scan(&friend)
		if err != nil {
			log.Println(err)
			return nil
		}
		friendid := fmt.Sprintf("%d", friend)
		val, err := infra.RDB.HGet(context.Background(), "user:username_map", friendid).Result()
		if err != nil {
			log.Println(err)
			return nil
		}
		friendUsername := models.FriendDTO{Friend: val}
		friends = append(friends, friendUsername)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("FRIENDS: ", friends)

	return friends
}

func (infra *Infra) FriendResponse(username string, action *models.FriendActionDTO) bool {
	val, err := infra.RDB.HMGet(context.Background(), "user:id_map", username, action.FriendID).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	receiver := val[0].(string)
	sender := val[1].(string)

	id, err := snowflake.ParseString(receiver)
	if err != nil {
		log.Println(err)
		return false
	}
	db := database.GetShardPool(infra.Pools, id)
	_, err = db.Exec(context.Background(), "UPDATE friends SET status=$1 WHERE user_id=$2 AND requester_id=$3", action.Action, receiver, sender)
	if err != nil {
		log.Println(err)
		return false
	}

	id, err = snowflake.ParseString(sender)
	if err != nil {
		log.Println(err)
		return false
	}
	db = database.GetShardPool(infra.Pools, id)
	_, err = db.Exec(context.Background(), "UPDATE friends SET status=$1 WHERE user_id=$2 AND requester_id=$3", action.Action, receiver, sender)
	if err != nil {
		log.Println(err)
		return false
	}
	key := fmt.Sprintf("conversation:%s:%s", username, action.FriendID)

	_, err = infra.RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: key,
		Values: map[string]interface{}{
			"username": "",
			"message":  "",
		},
	}).Result()
	if err != nil {
		log.Println(err)
	}

	return true
}

func (infra *Infra) GetLog(user1, user2 string) []models.Messages {
	/*key := fmt.Sprintf("conversation:%s:%s", user1, user2)
	res, err := infra.RDB.XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{key, "0"},
		Count:   100,
		Block:   300,
	}).Result()
	if err != nil {
		fmt.Println("ERROR")
	}*/
	key := fmt.Sprintf("conversation:%s:%s", user1, user2)
	res, err := infra.RDB.XRange(context.Background(), key, "-", "+").Result()
	if err != nil {
		log.Println("ERROR")
	}
	msgs := make([]models.Messages, 0, len(res))
	for _, e := range res {
		msgs = append(msgs, models.Messages{
			Username: e.Values["username"].(string),
			Message:  e.Values["message"].(string),
		})
		fmt.Println(e.Values["username"].(string))
		fmt.Println(e.Values["message"].(string))
	}
	return msgs
}
