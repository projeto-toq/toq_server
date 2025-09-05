package sessionservice

import "context"

// Service defines the public contract for session maintenance operations.
// Keep entrypoints minimal and observable (tracing/logging/metrics).
type Service interface {
	// CleanExpiredSessions deletes expired sessions within a transaction boundary.
	CleanExpiredSessions(ctx context.Context, limit int) (int64, error)
}
