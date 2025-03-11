package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"grpc-service-ref/internal/domain/models"
	"grpc-service-ref/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
    db *sql.DB
}

// New creates a new instance of the Postgres storage.
func New(connStr string) (*Storage, error) {
    const op = "storage.postgres.New"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
    const op = "storage.postgres.SaveUser"

    stmt, err :=  s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES($1, $2)")
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    res, err := stmt.ExecContext(ctx, email, passHash)
    if err != nil {
        var pgErr *pq.Error

        if errors.As(err, &pgErr) && pgErr.Code.Name() == "unique_violation" {
            return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
        }

        return 0, fmt.Errorf("%s: %w", op, err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
    const op = "storage.postgres.User"
    stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=$1")
    if err != nil {
        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    row := stmt.QueryRowContext(ctx, email)

    var user models.User
    err = row.Scan(&user.ID, &user.Email, &user.PassHash)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
        }

        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    return user, nil
}
