package postgres

import (
	"context"
	"fmt"
	"skillsrock-test-task/internal/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg config.PostgresConfig) (*Database, error) {
	ctx := context.Background()

	options := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%v?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	config, err := pgxpool.ParseConfig(options)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
