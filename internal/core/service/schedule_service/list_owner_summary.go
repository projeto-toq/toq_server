package scheduleservices

import (
	"context"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) ListOwnerSummary(ctx context.Context, filter schedulemodel.OwnerSummaryFilter) (schedulemodel.OwnerSummaryResult, error) {
	if filter.OwnerID <= 0 {
		return schedulemodel.OwnerSummaryResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return schedulemodel.OwnerSummaryResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.OwnerSummaryResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.list_owner_summary.tx_start_error", "err", txErr)
		return schedulemodel.OwnerSummaryResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.list_owner_summary.tx_rollback_error", "err", rbErr)
		}
	}()

	result, err := s.scheduleRepo.ListOwnerSummary(ctx, tx, filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_owner_summary.repo_error", "owner_id", filter.OwnerID, "err", err)
		return schedulemodel.OwnerSummaryResult{}, utils.InternalError("")
	}

	return result, nil
}

func validateRange(from, to time.Time) error {
	if from.IsZero() || to.IsZero() {
		return nil
	}
	if !from.Before(to) {
		return utils.ValidationError("range", "from must be before to")
	}
	return nil
}
