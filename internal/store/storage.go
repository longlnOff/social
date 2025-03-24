package store

import (
	"context"
	"database/sql"
)


type Storage struct {
	Post interface {
		Create(ctx context.Context,post *Post) error
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
