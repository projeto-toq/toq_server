package scheduleservices

import (
	"context"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) GetAvailability(ctx context.Context, filter schedulemodel.AvailabilityFilter) (AvailabilityResult, error) {
	if filter.ListingIdentityID <= 0 {
		return AvailabilityResult{}, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if filter.Range.From.IsZero() || filter.Range.To.IsZero() {
		return AvailabilityResult{}, utils.ValidationError("range", "from and to must be provided")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return AvailabilityResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return AvailabilityResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.get_availability.tx_start_error", "err", txErr)
		return AvailabilityResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.get_availability.tx_rollback_error", "err", rbErr)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if err != nil {
		return AvailabilityResult{}, mapAgendaError(ctx, logger, err, filter.ListingIdentityID)
	}

	loc, tzErr := resolveAgendaLocation(agenda)
	if tzErr != nil {
		return AvailabilityResult{}, tzErr
	}

	repoRange := schedulemodel.ScheduleRange{From: filter.Range.From, To: filter.Range.To, Loc: loc}
	repoFilter := buildAvailabilityRepoFilter(filter, repoRange, loc)

	data, dataErr := s.scheduleRepo.GetAvailabilityData(ctx, tx, repoFilter)
	if dataErr != nil {
		utils.SetSpanError(ctx, dataErr)
		logger.Error("schedule.get_availability.repo_error", "listing_identity_id", filter.ListingIdentityID, "err", dataErr)
		return AvailabilityResult{}, utils.InternalError("")
	}

	result := availabilityEngineCompute(repoRange, loc, filter.SlotDurationMinute, filter.Pagination, data)
	return result, nil
}
