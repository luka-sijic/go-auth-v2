package auto

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(Pools []*pgxpool.Pool) error {
	sql := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS friends (
			user_id BIGINT NOT NULL,
			requester_id BIGINT NOT NULL,
			STATUS TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			PRIMARY KEY (user_id, requester_id)
		);
	`

	for i, pool := range Pools {
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
