package scheduleservices

import (
	"context"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVisitEntry inserts an agenda entry linked to a visit.
// The entry is created as VISIT_PENDING or VISIT_CONFIRMED depending on the pending flag
// and is always blocking to avoid overlapping bookings.
func (s *scheduleService) CreateVisitEntry(ctx context.Context, agendaID uint64, visitID uint64, start, end time.Time, pending bool) (schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	entry := schedulemodel.NewAgendaEntry()
	entry.SetAgendaID(agendaID)
	entry.SetVisitID(visitID)
	entry.SetStartsAt(start)
	entry.SetEndsAt(end)
	entry.SetBlocking(true)
	if pending {
		entry.SetEntryType(schedulemodel.EntryTypeVisitPending)
	} else {
		entry.SetEntryType(schedulemodel.EntryTypeVisitConfirmed)
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return nil, txErr
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	id, insertErr := s.scheduleRepo.InsertEntry(ctx, tx, entry)
	if insertErr != nil {
		utils.SetSpanError(ctx, insertErr)
		return nil, insertErr
	}
	entry.SetID(id)

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		return nil, commitErr
	}
	committed = true

	return entry, nil
}
