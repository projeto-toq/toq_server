package scheduleservices

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteVisitEntry removes an agenda entry by its ID.
func (s *scheduleService) DeleteVisitEntry(ctx context.Context, entryID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return txErr
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if delErr := s.scheduleRepo.DeleteEntry(ctx, tx, entryID); delErr != nil {
		utils.SetSpanError(ctx, delErr)
		return delErr
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		return commitErr
	}
	committed = true

	return nil
}
