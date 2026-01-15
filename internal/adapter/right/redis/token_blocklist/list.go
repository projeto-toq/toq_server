package tokenblocklist

import (
	"context"
	"strings"
	"time"

	cacheport "github.com/projeto-toq/toq_server/internal/core/port/right/cache"
)

// List returns a coarse paginated slice of blocklisted JTIs using SCAN.
// Pagination is best-effort and not strongly ordered.
func (a *Adapter) List(ctx context.Context, page, pageSize int64) ([]cacheport.BlocklistItem, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}

	start := (page - 1) * pageSize
	cursor := uint64(0)
	items := make([]cacheport.BlocklistItem, 0, pageSize)
	skipped := int64(0)

	for {
		keys, next, err := a.client.Scan(ctx, cursor, keyPrefix+"*", pageSize*2).Result()
		if err != nil {
			return nil, err
		}

		for _, k := range keys {
			if skipped < start {
				skipped++
				continue
			}

			ttl, _ := a.client.TTL(ctx, k).Result()
			parts := strings.Split(k, ":")
			jti := parts[len(parts)-1]
			exp := time.Now().Add(ttl).Unix()
			items = append(items, cacheport.BlocklistItem{JTI: jti, ExpiresAt: exp})

			if int64(len(items)) >= pageSize {
				return items, nil
			}
		}

		if next == 0 {
			break
		}
		cursor = next
	}

	return items, nil
}
