package cache

import (
	"github.com/go-redis/redis/v8"
)

func NewValkeyClient(address string, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
