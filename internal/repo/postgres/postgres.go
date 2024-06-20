package postgres

import (
	"context"
	"fmt"

	"event-notification-bot/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name)

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{db: pool}, nil
}

func (p *Postgres) Close() {
	p.db.Close()
}
