package tokenblocklist

import (
	cacheport "github.com/projeto-toq/toq_server/internal/core/port/right/cache"
	"github.com/redis/go-redis/v9"
)

const keyPrefix = "toq:blocklist:jti:"

// Adapter implements TokenBlocklistPort using Redis keys with TTL.
type Adapter struct {
	client *redis.Client
}

// NewAdapter constructs a blocklist adapter backed by Redis.
func NewAdapter(client *redis.Client) cacheport.TokenBlocklistPort {
	return &Adapter{client: client}
}

func key(jti string) string {
	return keyPrefix + jti
}
