package sessionrepository

import (
	"context"
	"database/sql"
	"time"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

type SessionRepoPortInterface interface {
	CreateSession(ctx context.Context, tx *sql.Tx, session sessionmodel.SessionInterface) (err error)
	GetSessionByID(ctx context.Context, tx *sql.Tx, id int64) (session sessionmodel.SessionInterface, err error)
	GetActiveSessionByRefreshHash(ctx context.Context, tx *sql.Tx, hash string) (session sessionmodel.SessionInterface, err error)
	RevokeSession(ctx context.Context, tx *sql.Tx, id int64) error
	MarkSessionRotated(ctx context.Context, tx *sql.Tx, id int64) error
	RevokeSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error
	// DeleteSessionsByUserID permanently removes all sessions for a given user (any state)
	DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error
	GetActiveSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int64) (sessions []sessionmodel.SessionInterface, err error)
	UpdateSessionRotation(ctx context.Context, tx *sql.Tx, id int64, rotationCounter int, lastRefreshAt time.Time) error
	// DeleteExpiredSessions removes revoked+expired or absolutely expired sessions, returns affected count
	DeleteExpiredSessions(ctx context.Context, tx *sql.Tx, limit int) (int64, error)
}
