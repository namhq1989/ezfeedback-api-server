package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func getRedisConnFromURL(redisURL string) asynq.RedisClientOpt {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis url: %w", err))
	}

	return asynq.RedisClientOpt{
		Addr:     opt.Addr,
		Username: opt.Username,
		Password: opt.Password,
		DB:       0,
	}
}
