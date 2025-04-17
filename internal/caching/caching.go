package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Operations interface {
	GenerateKey(domain, key string) string

	Get(ctx *appcontext.AppContext, key string) (string, error)
	Set(ctx *appcontext.AppContext, key string, value interface{})
	SetTTL(ctx *appcontext.AppContext, key string, value interface{}, expiration time.Duration)
	Del(ctx *appcontext.AppContext, key string) (int64, error)
}

type Caching struct {
	redis *redis.Client
}

func NewCachingClient(redisURL string) *Caching {
	var (
		ctx    = context.Background()
		opt, _ = redis.ParseURL(redisURL)
	)

	// new client
	client := redis.NewClient(opt)

	// ping
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(client, redisotel.WithAttributes()); err != nil {
		panic(err)
	}

	fmt.Printf("⚡️ [caching]: connected \n")

	return &Caching{redis: client}
}
