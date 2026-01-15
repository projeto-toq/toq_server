package cache

import "context"

// TokenBlocklistPort defines operations to manage blocked access-token JTIs.
// Implementations should be lightweight and backed by an in-memory store like Redis.
type TokenBlocklistPort interface {
	Add(ctx context.Context, jti string, ttlSeconds int64) error
	Exists(ctx context.Context, jti string) (bool, error)
	Delete(ctx context.Context, jti string) error
	List(ctx context.Context, page, pageSize int64) (items []BlocklistItem, err error)
	Count(ctx context.Context) (int64, error)
}

// BlocklistItem represents a blocked JTI entry.
type BlocklistItem struct {
	JTI       string
	ExpiresAt int64 // unix seconds
}
