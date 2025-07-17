package database

import (
	"context"
	"hash/fnv"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	Pools []*pgxpool.Pool
	RDB   *redis.Client
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file", err)
	}
}

func Connect() {
	raw := os.Getenv("DATABASE_URLS")
	if raw == "" {
		log.Fatal("DATABASE_URLS is not set; should be comma-separated URLs")
	}
	urls := strings.Split(raw, ",")
	for i, u := range urls {
		pool, err := pgxpool.New(context.Background(), strings.TrimSpace(u))
		if err != nil {
			log.Fatalf("Unable to connect to database shard %d: %v", i, err)
		}
		log.Printf("Shard %d connected\n", i)
		Pools = append(Pools, pool)
	}
	ConnectRedis()
}

func ConnectRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	log.Println("Redis connected")
}

func Close() {
	for i, pool := range Pools {
		pool.Close()
		log.Printf("Shard %d closed\n", i)
	}
}

func GetShardPool(key snowflake.ID) *pgxpool.Pool {
	h := fnv.New32a()
	h.Write([]byte(key.String()))
	idx := int(h.Sum32()) % len(Pools)
	return Pools[idx]
}
