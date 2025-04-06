package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/longlnOff/social/internal/store"
)

var (
	USER_EXP_TIME = time.Duration(2 * time.Hour)
)

type UserStore struct {
	db *redis.Client
}

func NewUser(db *redis.Client) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (u *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	userIDKey := fmt.Sprintf("user:%d", userID)
	data, err := u.db.Get(ctx, userIDKey).Result()
	if err == redis.Nil {
		return nil, nil // Can't find in redis
	}
	if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (u *UserStore) Set(ctx context.Context, user *store.User) error {
	// Should check if user has ID first in production implmentation
	userIDKey := fmt.Sprintf("user:%d", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.db.SetEX(ctx, userIDKey, json, USER_EXP_TIME).Err()
}
