package photosessionservices

import (
	"context"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateServiceArea registers a new service area for the photographer.
func (s *photoSessionService) CreateServiceArea(ctx context.Context, input CreateServiceAreaInput) (ServiceAreaResult, error) {
	if input.PhotographerID == 0 {
		return ServiceAreaResult{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
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
		logger.Error("photo_session.service_area.create.tx_start_error", "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.service_area.create.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	existingAreas, err := s.repo.ListServiceAreasByPhotographer(ctx, tx, input.PhotographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.create.list_existing_error", "photographer_id", input.PhotographerID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if hasServiceAreaDuplicate(existingAreas, normalizedCity, normalizedState, 0) {
		return ServiceAreaResult{}, utils.ValidationError("city", "service area already exists for this city and state")
	}

	area := photosessionmodel.NewPhotographerServiceArea()
	area.SetPhotographerUserID(input.PhotographerID)
	area.SetCity(normalizedCity)
	area.SetState(normalizedState)

	areaID, err := s.repo.CreateServiceArea(ctx, tx, area)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.create.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	createdArea, err := s.repo.GetServiceAreaByID(ctx, tx, areaID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.create.get_created_error", "area_id", areaID, "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.service_area.create.tx_commit_error", "err", err)
		return ServiceAreaResult{}, utils.InternalError("")
	}
	committed = true

	return ServiceAreaResult{Area: createdArea}, nil
}
