package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/longlnOff/social/internal/store"
)

type Storage struct {
	User interface {
		Get(ctx context.Context, userID int64) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewCacheStorage(rdb *redis.Client) Storage {
	return Storage{
		User: NewUser(rdb),
	}
}
