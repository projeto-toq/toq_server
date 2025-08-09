package sessionrepository

import (
	context "context"
	"time"

	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"
)

type SessionRepoPortInterface interface {
	CreateSession(ctx context.Context, session sessionmodel.SessionInterface) (sessionmodel.SessionInterface, error)
	GetSessionByID(ctx context.Context, id int64) (sessionmodel.SessionInterface, error)
	GetActiveSessionByRefreshHash(ctx context.Context, hash string) (sessionmodel.SessionInterface, error)
	RevokeSession(ctx context.Context, id int64) error
	MarkSessionRotated(ctx context.Context, id int64) error
	RevokeSessionsByUserID(ctx context.Context, userID int64) error
	GetActiveSessionsByUserID(ctx context.Context, userID int64) ([]sessionmodel.SessionInterface, error)
	UpdateSessionRotation(ctx context.Context, id int64, rotationCounter int, lastRefreshAt time.Time) error
}
