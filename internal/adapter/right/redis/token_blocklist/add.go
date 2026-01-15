package tokenblocklist

import (
	"context"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// Add inserts a JTI into the blocklist with the given TTL in seconds.
func (a *Adapter) Add(ctx context.Context, jti string, ttlSeconds int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if ttlSeconds <= 0 {
		return fmt.Errorf("ttl must be positive")
	}

	return a.client.Set(ctx, key(jti), "1", time.Duration(ttlSeconds)*time.Second).Err()
}
