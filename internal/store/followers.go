package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)



type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func NewFollower(db *sql.DB) *FollowerStore {
	return &FollowerStore{
		db: db,
	}
}

func (s *FollowerStore) Follow(ctx context.Context, followerID int64, followedUserID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1,$2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, followedUserID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code.Name() {
            case "unique_violation":
                return ErrConflict
			default:
				return err
            }
        }
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID int64, followedUserID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, followedUserID, followerID)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}
