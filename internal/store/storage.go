package store

import (
	"context"
	"database/sql"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, string) (*User, error)
	}
	Sessions interface {
		Upsert(context.Context, *Session) error
		GetByUserID(context.Context, string) (*Session, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users:    &UserStore{db},
		Sessions: &SessionStore{db},
	}
}
