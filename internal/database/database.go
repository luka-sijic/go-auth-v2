package database

import (
	"context"
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
)

var shardDSNs = []string{
	"postgres://based@Staffer92@100.67.0.96/app",
}

func initShardPools(ctx context.Context, dsns []string) []*pgxpool.Pool {
	pools := make([]*pgxpool.Pool, len(dsns))
	for i, dsn := range dsns {
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Fatalf("failed to parse DSN for shard %d: %v", i, err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to connect to shard %d: %v", i, err)
		}
		pools[i] = pool
	}
	return pools
}
