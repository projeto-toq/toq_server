package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"errors"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAvailabilityData aggregates agenda rules and entries to support availability computation.
//
// Parameters:
//   - ctx: request-scoped context for tracing/logging.
//   - tx: optional transaction; use non-nil when composing with other repository calls.
//   - filter: AvailabilityFilter containing listingIdentityID and the desired time window.
//
// Returns: AvailabilityData with rules and entries; sql.ErrNoRows when the agenda does not exist; infra errors otherwise.
// Observability: tracer span, logger propagation, span error marking on infra failures; sql.ErrNoRows is propagated pure.
func (a *ScheduleAdapter) GetAvailabilityData(ctx context.Context, tx *sql.Tx, filter schedulemodel.AvailabilityFilter) (schedulemodel.AvailabilityData, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.AvailabilityData{}, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	agenda, err := a.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.schedule.availability.get_agenda_error", "listing_identity_id", filter.ListingIdentityID, "err", err)
		}
		return schedulemodel.AvailabilityData{}, err
	}

	entries, err := a.ListEntriesBetween(ctx, tx, agenda.ID(), filter.Range.From, filter.Range.To)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.availability.list_entries_error", "agenda_id", agenda.ID(), "err", err)
		return schedulemodel.AvailabilityData{}, err
	}

	rules, err := a.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.availability.list_rules_error", "agenda_id", agenda.ID(), "err", err)
		return schedulemodel.AvailabilityData{}, err
	}

	return schedulemodel.AvailabilityData{Entries: entries, Rules: rules}, nil
}
