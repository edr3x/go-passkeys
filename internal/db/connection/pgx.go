package connection

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnectPG(ctx context.Context) *pgxpool.Pool {
	pgconn := os.Getenv("POSTGRESQL_URL")
	pool, err := pgxpool.New(ctx, pgconn)
	if err != nil {
		log.Fatal(err)
	}

	pool.Config().MaxConns = 10
	pool.Config().MaxConnIdleTime = 20 * time.Second

	return pool
}
