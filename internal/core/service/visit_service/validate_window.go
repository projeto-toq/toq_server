package visitservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// validateWindow enforces lead time, horizon, and availability before creating a visit.
func (s *visitService) validateWindow(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface, input CreateVisitInput) error {
	if err := s.validateLeadTimeAndHorizon(input); err != nil {
		return err
	}

	return s.validateAvailability(ctx, tx, agenda, input)
}

func (s *visitService) validateLeadTimeAndHorizon(input CreateVisitInput) error {
	now := time.Now().UTC()
	leadLimit := now.Add(time.Duration(s.config.MinHoursAhead) * time.Hour)
	if !input.ScheduledStart.After(leadLimit) {
		return utils.ValidationError("scheduledStart", fmt.Sprintf("must be at least %d hours in advance", s.config.MinHoursAhead))
	}

	horizonLimit := now.Add(time.Duration(s.config.MaxDaysAhead) * 24 * time.Hour)
	if input.ScheduledStart.After(horizonLimit) || input.ScheduledEnd.After(horizonLimit) {
		return utils.ValidationError("scheduledStart", fmt.Sprintf("must be within %d days from now", s.config.MaxDaysAhead))
	}

	return nil
}

func (s *visitService) validateAvailability(ctx context.Context, _ *sql.Tx, agenda schedulemodel.AgendaInterface, input CreateVisitInput) error {
	logger := utils.LoggerFromContext(ctx)
	loc, tzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if tzErr != nil {
		return tzErr
	}

	duration := input.ScheduledEnd.Sub(input.ScheduledStart)
	if duration <= 0 {
		return utils.ValidationError("scheduledTime", "end must be after start")
	}

	if s.scheduleSvc == nil {
		logger.Error("visit.validate.availability.schedule_service_nil")
		return utils.InternalError("")
	}

	slot := schedulemodel.ScheduleRange{From: input.ScheduledStart, To: input.ScheduledEnd, Loc: loc}
	filter := schedulemodel.AvailabilityFilter{
		ListingIdentityID:  input.ListingIdentityID,
		Range:              slot,
		SlotDurationMinute: uint16(duration.Minutes()),
		Pagination:         schedulemodel.PaginationConfig{Page: 1, Limit: 1},
	}

	available, availErr := s.scheduleSvc.CheckSlotAvailability(ctx, filter, slot)
	if availErr != nil {
		return availErr
	}
	if !available {
		return utils.ConflictError("Requested slot is not available")
	}

	return nil
}
