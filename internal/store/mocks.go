package store

import (
	"context"
	"database/sql"
	"time"
)


func NewMockStore() Storage {
	return Storage{
		User: &MockUserStore{},

	}
}

type MockUserStore struct {}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (m *MockUserStore) GetByUserID(ctx context.Context, userID int64) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, token string, user *User, invitationExpirationDuration time.Duration) error {
	return nil
}

func (m *MockUserStore) ActivateUser(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
