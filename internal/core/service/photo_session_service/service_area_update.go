package photosessionservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateServiceArea updates the city and state of an existing service area.
func (s *photoSessionService) UpdateServiceArea(ctx context.Context, input UpdateServiceAreaInput) (ServiceAreaResult, error) {
	if input.PhotographerID == 0 {
		return ServiceAreaResult{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.ServiceAreaID == 0 {
		return ServiceAreaResult{}, utils.ValidationError("serviceAreaId", "serviceAreaId must be greater than zero")
	}

	normalizedCity, normalizedState, validationErr := normalizeServiceAreaFields(input.City, input.State)
	if validationErr != nil {
		return ServiceAreaResult{}, validationErr
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ServiceAreaResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.tx_start_error", "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.service_area.update.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	area, err := s.repo.GetServiceAreaByID(ctx, tx, input.ServiceAreaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceAreaResult{}, utils.NotFoundError("ServiceArea")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.get_error", "area_id", input.ServiceAreaID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if area.PhotographerUserID() != input.PhotographerID {
		return ServiceAreaResult{}, utils.AuthorizationError("service area does not belong to current user")
	}

	existingAreas, err := s.repo.ListServiceAreasByPhotographer(ctx, tx, input.PhotographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.list_existing_error", "photographer_id", input.PhotographerID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if hasServiceAreaDuplicate(existingAreas, normalizedCity, normalizedState, input.ServiceAreaID) {
		return ServiceAreaResult{}, utils.ValidationError("city", "service area already exists for this city and state")
	}

	area.SetCity(normalizedCity)
	area.SetState(normalizedState)

	if err := s.repo.UpdateServiceArea(ctx, tx, area); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceAreaResult{}, utils.NotFoundError("ServiceArea")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.repo_error", "area_id", input.ServiceAreaID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	updatedArea, err := s.repo.GetServiceAreaByID(ctx, tx, input.ServiceAreaID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.get_updated_error", "area_id", input.ServiceAreaID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.update.tx_commit_error", "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}
	committed = true

	return ServiceAreaResult{Area: updatedArea}, nil
}
