package scheduleservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CheckVisitConflict verifies if there is any blocking entry overlapping the requested interval.
// excludeEntryID allows ignoring a specific entry (useful when updating an existing visit entry).
func (s *scheduleService) CheckVisitConflict(ctx context.Context, agendaID uint64, start, end time.Time, excludeEntryID *uint64) (bool, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return false, err
	}
	defer spanEnd()

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return false, txErr
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	entries, listErr := s.scheduleRepo.ListEntriesBetween(ctx, tx, agendaID, start, end)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		return false, listErr
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		return false, commitErr
	}
	committed = true

	for _, e := range entries {
		if excludeEntryID != nil && e.ID() == *excludeEntryID {
			continue
		}
		if !e.Blocking() {
			continue
		}
		if e.StartsAt().Before(end) && e.EndsAt().After(start) {
			return true, nil
		}
	}

	return false, nil
}
