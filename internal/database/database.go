package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a pool: ", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: ", err)
	}

	return pool, nil
}
