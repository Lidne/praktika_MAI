package postgres

import (
	"context"
	"fmt"
	"github.com/Lidne/praktika_MAI/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	DBUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB)
	dbpool, err := pgxpool.New(ctx, DBUrl)
	if err != nil {
		log.Fatal("Error connecting to database")
		return nil, err
	}

	return dbpool, nil
}
