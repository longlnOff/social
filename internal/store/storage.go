package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resouce not found")
)


type Storage struct {
	Post interface {
		Create(ctx context.Context,post *Post) error
		Get(ctx context.Context, id int64) (*Post, error)
	}

	User interface {
		Create(ctx context.Context, user *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post: NewPost(db),
		User: NewUser(db),
	}
}
