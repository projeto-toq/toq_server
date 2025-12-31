package scheduleservices

import (
	"context"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateVisitEntryType changes the entry type (pending/confirmed) and blocking flag for a visit entry.
func (s *scheduleService) UpdateVisitEntryType(ctx context.Context, entryID uint64, newType schedulemodel.EntryType, blocking bool) (schedulemodel.AgendaEntryInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

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

	entry, getErr := s.scheduleRepo.GetEntryByID(ctx, tx, entryID)
	if getErr != nil {
		utils.SetSpanError(ctx, getErr)
		return nil, getErr
	}

	entry.SetEntryType(newType)
	entry.SetBlocking(blocking)

	if updateErr := s.scheduleRepo.UpdateEntry(ctx, tx, entry); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		return nil, updateErr
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		return nil, commitErr
	}
	committed = true

	return entry, nil
}
