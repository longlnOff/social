package store

import (
	"context"
	"database/sql"
)

type PostsStore struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *PostsStore {
	return &PostsStore{
		db: db,
	}
}

func (u *PostsStore) Create(ctx context.Context) error {
	return nil
}
