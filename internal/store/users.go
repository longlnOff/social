package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

func (u *UsersStore) Create(ctx context.Context) error {
	return nil
}
