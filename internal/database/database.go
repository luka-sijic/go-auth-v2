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

type App struct {
	Pools []*pgxpool.Pool
	Node  *snowflake.Node
	RDB   *redis.Client
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file", err)
	}
}

func NewApp() (*App, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Printf("snowflake node init failed: %v", err)
		return nil, err
	}
	pools := ConnectPostgres()
	rdb := ConnectRedis()

	return &App{Node: node, Pools: pools, RDB: rdb}, nil
}

func ConnectPostgres() []*pgxpool.Pool {
	raw := os.Getenv("DATABASE_URLS")
	if raw == "" {
		log.Fatal("DATABASE_URLS is not set; should be comma-separated URLs")
		return nil
	}
	urls := strings.Split(raw, ",")
	pools := make([]*pgxpool.Pool, 0, 4)
	for i, u := range urls {
		pool, err := pgxpool.New(context.Background(), strings.TrimSpace(u))
		if err != nil {
			log.Fatalf("Unable to connect to database shard %d: %v", i, err)
		}
		log.Printf("Shard %d connected\n", i)
		pools = append(pools, pool)
	}
	return pools
}

func ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	log.Println("Redis connected")
	return rdb
}

func (a *App) Close() {
	for i, pool := range a.Pools {
		pool.Close()
		log.Printf("Shard %d closed\n", i)
	}
}

func GetShardPool(pools []*pgxpool.Pool, key snowflake.ID) *pgxpool.Pool {
	//id := (key >> 12) & ((1 << 10) - 1)
	//fmt.Println("\033[32m POOL ID: ", id, " \033[0m")
	h := fnv.New32a()
	h.Write([]byte(key.String()))
	idx := int(h.Sum32()) % len(pools)
	return pools[idx]
}
