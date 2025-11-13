package mysqlscheduleadapter

import (
	"context"
	"database/sql"
	"fmt"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
)

func (a *ScheduleAdapter) GetAvailabilityData(ctx context.Context, tx *sql.Tx, filter schedulemodel.AvailabilityFilter) (schedulemodel.AvailabilityData, error) {
	agenda, err := a.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if err != nil {
		return schedulemodel.AvailabilityData{}, fmt.Errorf("get agenda for availability: %w", err)
	}

	entries, err := a.ListEntriesBetween(ctx, tx, agenda.ID(), filter.Range.From, filter.Range.To)
	if err != nil {
		return schedulemodel.AvailabilityData{}, fmt.Errorf("list entries for availability: %w", err)
	}

	rules, err := a.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		return schedulemodel.AvailabilityData{}, fmt.Errorf("list rules for availability: %w", err)
	}

	return schedulemodel.AvailabilityData{Entries: entries, Rules: rules}, nil
}
