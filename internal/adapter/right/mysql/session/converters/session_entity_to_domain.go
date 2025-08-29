package sessionconverters

import (
	"log/slog"
	"time"

	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func SessionEntityToDomain(entity []any) (session sessionmodel.SessionInterface, err error) {
	session = sessionmodel.NewSession()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting id to int64", "value", entity[0])
		return nil, utils.ErrInternalServer
	}
	session.SetID(id)

	userID, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting user_id to int64", "value", entity[1])
		return nil, utils.ErrInternalServer
	}
	session.SetUserID(userID)

	refreshHash, ok := entity[2].([]byte)
	if !ok {
		slog.Error("Error converting refresh_hash to []byte", "value", entity[2])
		return nil, utils.ErrInternalServer
	}
	session.SetRefreshHash(string(refreshHash))

	if entity[3] != nil {
		tokenJTI, ok := entity[3].([]byte)
		if !ok {
			slog.Error("Error converting token_jti to []byte", "value", entity[3])
			return nil, utils.ErrInternalServer
		}
		session.SetTokenJTI(string(tokenJTI))
	}

	expiresAt, ok := entity[4].(time.Time)
	if !ok {
		slog.Error("Error converting expires_at to time.Time", "value", entity[4])
		return nil, utils.ErrInternalServer
	}
	session.SetExpiresAt(expiresAt)

	if entity[5] != nil {
		absoluteExpiresAt, ok := entity[5].(time.Time)
		if !ok {
			slog.Error("Error converting absolute_expires_at to time.Time", "value", entity[5])
			return nil, utils.ErrInternalServer
		}
		session.SetAbsoluteExpiresAt(absoluteExpiresAt)
	}

	createdAt, ok := entity[6].(time.Time)
	if !ok {
		slog.Error("Error converting created_at to time.Time", "value", entity[6])
		return nil, utils.ErrInternalServer
	}
	session.SetCreatedAt(createdAt)

	if entity[7] != nil {
		rotatedAt, ok := entity[7].(time.Time)
		if !ok {
			slog.Error("Error converting rotated_at to time.Time", "value", entity[7])
			return nil, utils.ErrInternalServer
		}
		session.SetRotatedAt(&rotatedAt)
	}

	if entity[8] != nil {
		userAgent, ok := entity[8].([]byte)
		if !ok {
			slog.Error("Error converting user_agent to []byte", "value", entity[8])
			return nil, utils.ErrInternalServer
		}
		session.SetUserAgent(string(userAgent))
	}

	if entity[9] != nil {
		ip, ok := entity[9].([]byte)
		if !ok {
			slog.Error("Error converting ip to []byte", "value", entity[9])
			return nil, utils.ErrInternalServer
		}
		session.SetIP(string(ip))
	}

	if entity[10] != nil {
		deviceID, ok := entity[10].([]byte)
		if !ok {
			slog.Error("Error converting device_id to []byte", "value", entity[10])
			return nil, utils.ErrInternalServer
		}
		session.SetDeviceID(string(deviceID))
	}

	if entity[11] != nil {
		rotationCounter, ok := entity[11].(int64)
		if !ok {
			slog.Error("Error converting rotation_counter to int64", "value", entity[11])
			return nil, utils.ErrInternalServer
		}
		session.SetRotationCounter(int(rotationCounter))
	}

	if entity[12] != nil {
		lastRefreshAt, ok := entity[12].(time.Time)
		if !ok {
			slog.Error("Error converting last_refresh_at to time.Time", "value", entity[12])
			return nil, utils.ErrInternalServer
		}
		session.SetLastRefreshAt(&lastRefreshAt)
	}

	revoked, ok := entity[13].(int64)
	if !ok {
		slog.Error("Error converting revoked to bool", "value", entity[13])
		return nil, utils.ErrInternalServer
	}
	session.SetRevoked(revoked == 1)

	return session, nil
}
