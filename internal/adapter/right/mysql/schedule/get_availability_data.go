package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetAvailabilityData aggregates entries and rules for a listing to compute availability windows.
// Returns sql.ErrNoRows when the agenda does not exist; other infra errors are marked on the span and logged.
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
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.availability.get_agenda_error", "listing_identity_id", filter.ListingIdentityID, "err", err)
		return schedulemodel.AvailabilityData{}, fmt.Errorf("get agenda for availability: %w", err)
	}

	entries, err := a.ListEntriesBetween(ctx, tx, agenda.ID(), filter.Range.From, filter.Range.To)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.availability.list_entries_error", "agenda_id", agenda.ID(), "err", err)
		return schedulemodel.AvailabilityData{}, fmt.Errorf("list entries for availability: %w", err)
	}

	rules, err := a.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.schedule.availability.list_rules_error", "agenda_id", agenda.ID(), "err", err)
		return schedulemodel.AvailabilityData{}, fmt.Errorf("list rules for availability: %w", err)
	}

	return schedulemodel.AvailabilityData{Entries: entries, Rules: rules}, nil
}
