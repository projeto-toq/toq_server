package photosessionservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetServiceArea retrieves a single service area owned by the photographer.
func (s *photoSessionService) GetServiceArea(ctx context.Context, input ServiceAreaDetailInput) (ServiceAreaResult, error) {
	if input.PhotographerID == 0 {
		return ServiceAreaResult{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.ServiceAreaID == 0 {
		return ServiceAreaResult{}, utils.ValidationError("serviceAreaId", "serviceAreaId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.GetServiceArea")
	if err != nil {
		return ServiceAreaResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.detail.tx_start_error", "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}
	defer func() {
		if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
			utils.SetSpanError(ctx, rollbackErr)
			logger.Error("photo_session.service_area.detail.tx_rollback_error", "err", rollbackErr)
		}
	}()

	area, err := s.repo.GetServiceAreaByID(ctx, tx, input.ServiceAreaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceAreaResult{}, utils.NotFoundError("ServiceArea")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.detail.repo_error", "area_id", input.ServiceAreaID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if area.PhotographerUserID() != input.PhotographerID {
		return ServiceAreaResult{}, utils.AuthorizationError("service area does not belong to current user")
	}

	return ServiceAreaResult{Area: area}, nil
}
