package sessionconverters

import (
	"context"

	sessionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/session/entities"
	sessionmodel "github.com/projeto-toq/toq_server/internal/core/model/session_model"
)

func SessionDomainToEntity(ctx context.Context, session sessionmodel.SessionInterface) sessionentities.SessionEntity {
	return sessionentities.SessionEntity{
		ID:                session.GetID(),
		UserID:            session.GetUserID(),
		RefreshHash:       session.GetRefreshHash(),
		TokenJTI:          session.GetTokenJTI(),
		ExpiresAt:         session.GetExpiresAt(),
		AbsoluteExpiresAt: session.GetAbsoluteExpiresAt(),
		CreatedAt:         session.GetCreatedAt(),
		RotatedAt:         session.GetRotatedAt(),
		UserAgent:         session.GetUserAgent(),
		IP:                session.GetIP(),
		DeviceID:          session.GetDeviceID(),
		RotationCounter:   session.GetRotationCounter(),
		LastRefreshAt:     session.GetLastRefreshAt(),
		Revoked:           session.GetRevoked(),
	}
}
