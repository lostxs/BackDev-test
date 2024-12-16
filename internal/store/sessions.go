package store

import (
	"context"
	"database/sql"
	"errors"
)

type Session struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	RefreshTokenHash string `json:"refresh_token_hash"`
}

type SessionStore struct {
	db *sql.DB
}

var (
	ErrSessionNotFound = errors.New("session not found")
)

func (s *SessionStore) Upsert(ctx context.Context, session *Session) error {
	query := `
	INSERT INTO sessions (user_id, refresh_token_hash) 
	VALUES ($1, $2) 
	ON CONFLICT (user_id) DO UPDATE SET refresh_token_hash = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		session.UserID,
		session.RefreshTokenHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionStore) GetByUserID(ctx context.Context, userID string) (*Session, error) {
	query := `
	SELECT id, user_id, refresh_token_hash 
	FROM sessions 
	WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var session Session
	row := s.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}
