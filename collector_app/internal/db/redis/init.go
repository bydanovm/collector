package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewUserRedis(opts ...rdsOptions) UserRedisImpl {
	options := &redis.Options{
		Addr:     "redis:6379",
		Password: "1234567890",
		DB:       0,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return UserRedisImpl{
		Redis:   client,
		process: &Process{},
	}
}
