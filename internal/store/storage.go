package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resouce not found")
	ErrConflict = errors.New("resouce already exists")
	QueryTimeoutDuration = 5 * time.Second
)


type Storage struct {
	Post interface {
		Create(ctx context.Context,post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		GetUserFeed(ctx context.Context, userID int64, p PaginatedFeed) ([]PostWithMetadata, error)
	}

	User interface {
		Create(ctx context.Context, user *User) error
		GetByUserID(ctx context.Context, userID int64) (*User, error)
	}

	Comment interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}

	Follower interface {
		Follow(ctx context.Context, followerID int64, followedUserID int64) error
		Unfollow(ctx context.Context, followerID int64, followedUserID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post: 		NewPost(db),
		User: 		NewUser(db),
		Comment: 	NewComment(db),
		Follower: 	NewFollower(db),
	}
}
