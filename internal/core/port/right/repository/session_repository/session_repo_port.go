package sessionrepository

import (
	"context"
	"database/sql"
	"time"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

// SessionRepoPortInterface defines persistence operations for authentication sessions.
//
// Contract highlights:
//   - Methods accept optional transactions (tx may be nil for maintenance reads/writes when allowed).
//   - Implementations must use InstrumentedAdapter helpers for tracing/metrics/logging.
//   - Return sql.ErrNoRows where applicable; do not wrap domain/HTTP concerns at this layer.
type SessionRepoPortInterface interface {
	// CreateSession inserts a new session record and populates the ID on success.
	// Returns infrastructure errors (constraint violations, connectivity) as-is.
	CreateSession(ctx context.Context, tx *sql.Tx, session sessionmodel.SessionInterface) error

	// GetSessionByID fetches a session by primary key; returns sql.ErrNoRows when absent.
	GetSessionByID(ctx context.Context, tx *sql.Tx, id int64) (session sessionmodel.SessionInterface, err error)

	// GetActiveSessionByRefreshHash retrieves a non-revoked, non-expired session by refresh_hash; sql.ErrNoRows if not found/expired.
	GetActiveSessionByRefreshHash(ctx context.Context, tx *sql.Tx, hash string) (session sessionmodel.SessionInterface, err error)

	// RevokeSession marks a session as revoked by ID.
	RevokeSession(ctx context.Context, tx *sql.Tx, id int64) error

	// MarkSessionRotated sets rotated_at to now for a given session ID.
	MarkSessionRotated(ctx context.Context, tx *sql.Tx, id int64) error

	// RevokeSessionsByUserID revokes all active sessions of a user.
	RevokeSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error

	// DeleteSessionsByUserID permanently removes all sessions for a given user (any state).
	DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error

	// GetActiveSessionsByUserID lists active (non-revoked, non-expired) sessions for a user; returns empty slice when none.
	GetActiveSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) (sessions []sessionmodel.SessionInterface, err error)

	// UpdateSessionRotation updates rotation counter and last_refresh_at for a session.
	UpdateSessionRotation(ctx context.Context, tx *sql.Tx, id int64, rotationCounter int, lastRefreshAt time.Time) error

	// DeleteExpiredSessions removes any session whose expires_at is past now or whose absolute_expires_at is past now.
	// Intended for maintenance jobs; returns affected rows count.
	DeleteExpiredSessions(ctx context.Context, tx *sql.Tx, limit int) (int64, error)

	// DeleteExpiredSessionsBefore removes sessions whose effective expiry is before cutoff (absolute or sliding), capped by limit.
	DeleteExpiredSessionsBefore(ctx context.Context, tx *sql.Tx, cutoff time.Time, limit int) (int64, error)
}
