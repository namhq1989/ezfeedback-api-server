package caching

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
)

const (
	defaultTTL = time.Hour * 24 * 30
)

func (c Caching) Set(ctx *appcontext.AppContext, key string, value interface{}) {
	// default ttl of this function is 30 days
	c.SetTTL(ctx, key, value, defaultTTL)
}

func (c Caching) SetTTL(ctx *appcontext.AppContext, key string, value interface{}, expiration time.Duration) {
	var dataToStore []byte

	// handle string
	if strVal, ok := value.(string); ok {
		dataToStore = []byte(strVal)
	} else {
		dataToStore, _ = json.Marshal(value)
	}

	c.redis.Set(ctx.Context(), key, dataToStore, expiration)
}

func (c Caching) Get(ctx *appcontext.AppContext, key string) (string, error) {
	data, err := c.redis.Get(ctx.Context(), key).Result()
	if err == nil {
		return data, nil
	}
	if err.Error() == "redis: nil" {
		return "", nil
	}
	return "", err
}

func (c Caching) Del(ctx *appcontext.AppContext, key string) (int64, error) {
	return c.redis.Del(ctx.Context(), key).Result()
}

func (Caching) GenerateKey(domain, key string) string {
	return fmt.Sprintf("ezfb_apiServer:caching:%s:%s", domain, key)
}
