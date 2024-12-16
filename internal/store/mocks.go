package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type MockUserStore struct {
	users map[string]*User
}

type MockSessionStore struct {
	sessions map[string]*Session
}

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{
			users: make(map[string]*User),
		},
		Sessions: &MockSessionStore{
			sessions: make(map[string]*Session),
		},
	}
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	if _, exists := m.users[user.Email]; exists {
		return ErrDuplicateEmail
	}

	m.users[user.ID] = user
	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, id string) (*User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidUserID
	}

	return nil, ErrUserNotFound
}

func (m *MockSessionStore) Upsert(ctx context.Context, session *Session) error {
	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionStore) GetByUserID(ctx context.Context, userID string) (*Session, error) {
	if session, exists := m.sessions[userID]; exists {
		return session, nil
	}
	return nil, ErrSessionNotFound
}

func (m *MockSessionStore) Delete(ctx context.Context, id string) error {
	delete(m.sessions, id)
	return nil
}
