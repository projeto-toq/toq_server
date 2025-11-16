package mediaprocessingservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// resolveRequestedBy ensures every operation carries the user responsible for the change.
func (s *mediaProcessingService) resolveRequestedBy(ctx context.Context, requestedBy uint64) (uint64, error) {
	if requestedBy > 0 {
		return requestedBy, nil
	}

	userID, err := s.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return 0, derrors.Auth("authentication required")
	}
	if userID <= 0 {
		return 0, derrors.Auth("authentication required")
	}
	return uint64(userID), nil
}
