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

	loc := utils.DetermineRangeLocation(filter.Range.From, filter.Range.To, nil)
	repoFilter := filter
	repoFilter.Range.From, repoFilter.Range.To = utils.NormalizeRangeToUTC(filter.Range.From, filter.Range.To, loc)
	repoFilter.Range.Loc = time.UTC

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

	result, err := s.scheduleRepo.ListOwnerSummary(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_owner_summary.repo_error", "owner_id", filter.OwnerID, "err", err)
		return schedulemodel.OwnerSummaryResult{}, utils.InternalError("")
	}

	for i, item := range result.Items {
		entries := make([]schedulemodel.SummaryEntry, 0, len(item.Entries))
		for _, entry := range item.Entries {
			entries = append(entries, schedulemodel.SummaryEntry{
				EntryType: entry.EntryType,
				StartsAt:  utils.ConvertToLocation(entry.StartsAt, loc),
				EndsAt:    utils.ConvertToLocation(entry.EndsAt, loc),
				Blocking:  entry.Blocking,
			})
		}
		result.Items[i].Entries = entries
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
