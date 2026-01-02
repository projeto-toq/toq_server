package scheduleservices

import (
	"context"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CheckSlotAvailability verifies if the requested slot fits the agenda using the agenda timezone.
func (s *scheduleService) CheckSlotAvailability(ctx context.Context, filter schedulemodel.AvailabilityFilter, slot schedulemodel.ScheduleRange) (bool, error) {
	if filter.ListingIdentityID <= 0 {
		return false, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if slot.From.IsZero() || slot.To.IsZero() {
		return false, utils.ValidationError("range", "from and to must be provided")
	}
	if err := validateRange(slot.From, slot.To); err != nil {
		return false, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.check_slot_availability.tx_start_error", "err", txErr)
		return false, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.check_slot_availability.tx_rollback_error", "err", rbErr)
		}
	}()

	agenda, agErr := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if agErr != nil {
		return false, mapAgendaError(ctx, logger, agErr, filter.ListingIdentityID)
	}

	loc, tzErr := resolveAgendaLocation(agenda)
	if tzErr != nil {
		return false, tzErr
	}

	normalizedSlot := schedulemodel.ScheduleRange{From: utils.ConvertToLocation(slot.From, loc), To: utils.ConvertToLocation(slot.To, loc), Loc: loc}
	repoFilter := buildAvailabilityRepoFilter(filter, normalizedSlot, loc)

	data, dataErr := s.scheduleRepo.GetAvailabilityData(ctx, tx, repoFilter)
	if dataErr != nil {
		utils.SetSpanError(ctx, dataErr)
		logger.Error("schedule.check_slot_availability.repo_error", "listing_identity_id", filter.ListingIdentityID, "err", dataErr)
		return false, utils.InternalError("")
	}

	fits := availabilityEngineFits(normalizedSlot, loc, data)
	return fits, nil
}
