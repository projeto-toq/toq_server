package sessionconverters

import (
	"fmt"
	"time"

	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

func SessionEntityToDomain(entity []any) (session sessionmodel.SessionInterface, err error) {
	session = sessionmodel.NewSession()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("convert id to int64: %T", entity[0])
	}
	session.SetID(id)

	userID, ok := entity[1].(int64)
	if !ok {
		return nil, fmt.Errorf("convert user_id to int64: %T", entity[1])
	}
	session.SetUserID(userID)

	refreshHash, ok := entity[2].([]byte)
	if !ok {
		return nil, fmt.Errorf("convert refresh_hash to []byte: %T", entity[2])
	}
	session.SetRefreshHash(string(refreshHash))

	if entity[3] != nil {
		tokenJTI, ok := entity[3].([]byte)
		if !ok {
			return nil, fmt.Errorf("convert token_jti to []byte: %T", entity[3])
		}
		session.SetTokenJTI(string(tokenJTI))
	}

	expiresAt, ok := entity[4].(time.Time)
	if !ok {
		return nil, fmt.Errorf("convert expires_at to time.Time: %T", entity[4])
	}
	session.SetExpiresAt(expiresAt)

	if entity[5] != nil {
		absoluteExpiresAt, ok := entity[5].(time.Time)
		if !ok {
			return nil, fmt.Errorf("convert absolute_expires_at to time.Time: %T", entity[5])
		}
		session.SetAbsoluteExpiresAt(absoluteExpiresAt)
	}

	createdAt, ok := entity[6].(time.Time)
	if !ok {
		return nil, fmt.Errorf("convert created_at to time.Time: %T", entity[6])
	}
	session.SetCreatedAt(createdAt)

	if entity[7] != nil {
		rotatedAt, ok := entity[7].(time.Time)
		if !ok {
			return nil, fmt.Errorf("convert rotated_at to time.Time: %T", entity[7])
		}
		session.SetRotatedAt(&rotatedAt)
	}

	if entity[8] != nil {
		userAgent, ok := entity[8].([]byte)
		if !ok {
			return nil, fmt.Errorf("convert user_agent to []byte: %T", entity[8])
		}
		session.SetUserAgent(string(userAgent))
	}

	if entity[9] != nil {
		ip, ok := entity[9].([]byte)
		if !ok {
			return nil, fmt.Errorf("convert ip to []byte: %T", entity[9])
		}
		session.SetIP(string(ip))
	}

	if entity[10] != nil {
		deviceID, ok := entity[10].([]byte)
		if !ok {
			return nil, fmt.Errorf("convert device_id to []byte: %T", entity[10])
		}
		session.SetDeviceID(string(deviceID))
	}

	if entity[11] != nil {
		rotationCounter, ok := entity[11].(int64)
		if !ok {
			return nil, fmt.Errorf("convert rotation_counter to int64: %T", entity[11])
		}
		session.SetRotationCounter(int(rotationCounter))
	}

	if entity[12] != nil {
		lastRefreshAt, ok := entity[12].(time.Time)
		if !ok {
			return nil, fmt.Errorf("convert last_refresh_at to time.Time: %T", entity[12])
		}
		session.SetLastRefreshAt(&lastRefreshAt)
	}

	revoked, ok := entity[13].(int64)
	if !ok {
		return nil, fmt.Errorf("convert revoked to int64: %T", entity[13])
	}
	session.SetRevoked(revoked == 1)

	return session, nil
}
