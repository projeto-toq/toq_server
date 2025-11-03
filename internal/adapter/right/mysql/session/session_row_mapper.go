package sessionmysqladapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func (sa *SessionAdapter) mapSessionFromScanner(ctx context.Context, scanner rowScanner, operation string) (sessionmodel.SessionInterface, error) {
	var (
		id                int64
		userID            int64
		refreshHash       []byte
		tokenJTI          sql.NullString
		expiresAt         time.Time
		absoluteExpiresAt sql.NullTime
		createdAt         time.Time
		rotatedAt         sql.NullTime
		userAgent         sql.NullString
		ip                sql.NullString
		deviceID          sql.NullString
		rotationCounter   sql.NullInt64
		lastRefreshAt     sql.NullTime
		revoked           bool
	)

	if err := scanner.Scan(&id, &userID, &refreshHash, &tokenJTI, &expiresAt, &absoluteExpiresAt, &createdAt, &rotatedAt, &userAgent, &ip, &deviceID, &rotationCounter, &lastRefreshAt, &revoked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error(fmt.Sprintf("mysql.session.%s.scan_error", operation), "error", err)
		return nil, fmt.Errorf("scan session (%s): %w", operation, err)
	}

	session := sessionmodel.NewSession()
	session.SetID(id)
	session.SetUserID(userID)
	session.SetRefreshHash(string(refreshHash))
	if tokenJTI.Valid {
		session.SetTokenJTI(tokenJTI.String)
	}
	session.SetExpiresAt(expiresAt)
	if absoluteExpiresAt.Valid {
		session.SetAbsoluteExpiresAt(absoluteExpiresAt.Time)
	}
	session.SetCreatedAt(createdAt)
	if rotatedAt.Valid {
		rotated := rotatedAt.Time
		session.SetRotatedAt(&rotated)
	}
	if userAgent.Valid {
		session.SetUserAgent(userAgent.String)
	}
	if ip.Valid {
		session.SetIP(ip.String)
	}
	if deviceID.Valid {
		session.SetDeviceID(deviceID.String)
	}
	if rotationCounter.Valid {
		session.SetRotationCounter(int(rotationCounter.Int64))
	}
	if lastRefreshAt.Valid {
		last := lastRefreshAt.Time
		session.SetLastRefreshAt(&last)
	}
	session.SetRevoked(revoked)

	return session, nil
}
