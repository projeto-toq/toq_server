package photosessionservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteServiceArea removes a service area owned by the photographer.
func (s *photoSessionService) DeleteServiceArea(ctx context.Context, input DeleteServiceAreaInput) error {
	if input.PhotographerID == 0 {
		return utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.ServiceAreaID == 0 {
		return utils.ValidationError("serviceAreaId", "serviceAreaId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.delete.tx_start_error", "err", err)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.service_area.delete.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	area, err := s.repo.GetServiceAreaByID(ctx, tx, input.ServiceAreaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("ServiceArea")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.delete.get_error", "area_id", input.ServiceAreaID, "err", err)
		return utils.InternalError("")
	}

	if area.PhotographerUserID() != input.PhotographerID {
		return utils.AuthorizationError("service area does not belong to current user")
	}

	if err := s.repo.DeleteServiceArea(ctx, tx, input.ServiceAreaID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("ServiceArea")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.delete.repo_error", "area_id", input.ServiceAreaID, "err", err)
		return utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.delete.tx_commit_error", "err", err)
		return utils.InternalError("")
	}
	committed = true

	return nil
}
