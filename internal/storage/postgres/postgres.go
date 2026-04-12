package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rudrapatel51/goproject1-student/internal/config"
	"github.com/rudrapatel51/goproject1-student/internal/storage"
	"github.com/rudrapatel51/goproject1-student/internal/types"
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


func (s *Storage) GetStudentByID(id int) (types.Student, error) {
	const query = `
SELECT id, name, email, age, created_at
FROM students
WHERE id = $1;`

	var student types.Student
	err := s.db.QueryRow(context.Background(), query, id).Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.Student{}, storage.ErrStudentNotFound
		}
		return types.Student{}, fmt.Errorf("query student by id: %w", err)
	}

	return student, nil
}

func (s *Storage) GetAllStudents() ([]types.Student, error) {
	const query = `
SELECT id, name, email, age, created_at
FROM students
ORDER BY id;`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query all students: %w", err)
	}
	defer rows.Close()

	students := make([]types.Student, 0)
	for rows.Next() {
		var student types.Student
		if err = rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan student row: %w", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student rows: %w", err)
	}

	return students, nil
}


func (s *Storage) UpdateById(id int, name string, email string, age int) error {
	const query = `
UPDATE students
SET name = $1, email = $2, age = $3
WHERE id = $4;`
	result, err := s.db.Exec(context.Background(), query, name, email, age, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "students_email_key" {
			return storage.ErrStudentEmailAlreadyExists
		}
		return fmt.Errorf("update student: %w", err)
	}
	if result.RowsAffected() == 0 {
		return storage.ErrStudentNotFound
	}
	return nil
}

func (s *Storage) DeleteById(id int) error {
	const query = `
DELETE FROM students
WHERE id = $1;`
	result, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("delete student: %w", err)
	}
	if result.RowsAffected() == 0 {
		return storage.ErrStudentNotFound
	}
	return nil
}