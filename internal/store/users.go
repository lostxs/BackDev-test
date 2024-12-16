package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrInvalidUserID  = errors.New("invalid user id")
	ErrUserNotFound   = errors.New("user not found")
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (email) 
	VALUES ($1) 
	RETURNING id, email
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Email,
	).Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	query := `
	SELECT id, email 
	FROM users 
	WHERE id = $1 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err = s.db.QueryRowContext(
		ctx,
		query,
		idUUID,
	).Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
