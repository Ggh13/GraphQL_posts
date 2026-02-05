package postgres

import (
	"context"
	"fmt"
	"qraphQL_posts/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost" yaml:"POSTGRES_HOST"`
	Port     uint16 `env:"POSTGRES_PORT" env-default:"5432"      yaml:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER" env-default:"root"      yaml:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASS" env-default:"1234"      yaml:"POSTGRES_PASS"`
	Name     string `env:"POSTGRES_DB"   env-default:"postgres"  yaml:"POSTGRES_DB"`

	MinConns int32 `env:"POSTGRES_MIN_CONN" env-default:"2"  yaml:"POSTGRES_MIN_CONN"`
	MaxConns int32 `env:"POSTGRES_MAX_CONN" env-default:"10" yaml:"POSTGRES_MAX_CONN"`
}

func NewPostgres(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.MaxConns,
		cfg.MinConns,
	)

	logger.GetLoggerFromCtx(ctx).Info(ctx, connString)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Создаем таблицы если их нет
	createTablesSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			surname VARCHAR(255) NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id INT NOT NULL,
			comments_disabled BOOLEAN DEFAULT FALSE
		);
		
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL,
			parent_id INT,
			author_id INT NOT NULL,
			content VARCHAR(2000) NOT NULL
		);
	`

	_, err = conn.Exec(ctx, createTablesSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Database tables created/verified")

	return conn, nil
}
