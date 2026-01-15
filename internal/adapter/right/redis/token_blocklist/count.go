package tokenblocklist

import "context"

// Count returns number of blocklisted JTIs (best-effort via SCAN).
func (a *Adapter) Count(ctx context.Context) (int64, error) {
	var cursor uint64
	var total int64

	for {
		keys, next, err := a.client.Scan(ctx, cursor, keyPrefix+"*", 1000).Result()
		if err != nil {
			return 0, err
		}
		total += int64(len(keys))
		if next == 0 {
			break
		}
		cursor = next
	}

	return total, nil
}
