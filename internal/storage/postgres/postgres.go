package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudrapatel51/goproject1-student/internal/config"
	"github.com/rudrapatel51/goproject1-student/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}
	poolCfg.MaxConns = cfg.Postgres.MaxConns

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	const schema = `
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    age INT NOT NULL CHECK (age >= 0 AND age <= 150),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

	if _, err = db.Exec(ctx, schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("ensure students table: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CreateStudent(name string, email string, age int) (int64, error) {
	const query = `
INSERT INTO students (name, email, age)
VALUES ($1, $2, $3)
RETURNING id;`

	var id int64
	if err := s.db.QueryRow(context.Background(), query, name, email, age).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "students_email_key" {
			return 0, storage.ErrStudentEmailAlreadyExists
		}
		return 0, fmt.Errorf("insert student: %w", err)
	}

	return id, nil
}
