package auto

import (
	"app/internal/database"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init() error {
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`

	for i, pool := range database.Pools {
		if err := execOnShard(pool, sql); err != nil {
			return fmt.Errorf("shard %d: %w", i, err)
		}
	}
	return nil
}

func execOnShard(pool *pgxpool.Pool, sql string) error {
	_, err := pool.Exec(context.Background(), sql)
	return err
}
