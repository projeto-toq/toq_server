package goroutines

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"
)

// SessionCleaner periodically deletes expired or fully expired sessions.
func SessionCleaner(db *sql.DB, wg *sync.WaitGroup, ctx context.Context, interval time.Duration) {
	slog.Info("session cleaner routine started")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Info("session cleaner routine stopped")
			return
		case <-ticker.C:
			_, err := db.Exec(`DELETE FROM sessions WHERE (expires_at < UTC_TIMESTAMP() AND revoked = 1) OR (absolute_expires_at IS NOT NULL AND absolute_expires_at < UTC_TIMESTAMP()) LIMIT 500`)
			if err != nil {
				slog.Warn("session cleaner delete failed", "err", err)
			}
		}
	}
}
