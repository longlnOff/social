package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resouce not found")
	QueryTimeoutDuration = 5 * time.Second
)


type Storage struct {
	Post interface {
		Create(ctx context.Context,post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
	}

	User interface {
		Create(ctx context.Context, user *User) error
	}

	Comment interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post: 		NewPost(db),
		User: 		NewUser(db),
		Comment: 	NewComment(db),
	}
}
