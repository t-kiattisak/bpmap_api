package redis

import (
	"context"
	"fmt"
	"pbmap_api/src/pkg/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return client, nil
}
