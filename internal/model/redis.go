package model

import (
	"UserService/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func ConnectRedis(ctx context.Context, cfg *config.RedisDatabase) *redis.Client {
	cl := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	err := cl.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}

	return cl
}
