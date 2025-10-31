package photosessionservices

import (
	"context"
	"sort"
	"time"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListTimeOff returns paginated agenda entries of type TIME_OFF.
func (s *photoSessionService) ListTimeOff(ctx context.Context, input ListTimeOffInput) (ListTimeOffOutput, error) {
	if input.PhotographerID == 0 {
		return ListTimeOffOutput{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.RangeFrom.IsZero() {
		return ListTimeOffOutput{}, utils.ValidationError("rangeFrom", "rangeFrom is required")
	}
	if input.RangeTo.IsZero() {
		return ListTimeOffOutput{}, utils.ValidationError("rangeTo", "rangeTo is required")
	}
	if input.RangeTo.Before(input.RangeFrom) {
		return ListTimeOffOutput{}, utils.ValidationError("rangeTo", "rangeTo must be greater than or equal to rangeFrom")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListTimeOff")
	if err != nil {
		return ListTimeOffOutput{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc := input.Location
	if loc == nil {
		loc = time.UTC
	}

	fromLocal := utils.ConvertToLocation(input.RangeFrom, loc)
	toLocal := utils.ConvertToLocation(input.RangeTo, loc)

	page := input.Page
	if page <= 0 {
		page = defaultAgendaPage
	}

	size := input.Size
	if size <= 0 {
		size = defaultAgendaSize
	}
	if size > maxAgendaPageSize {
		size = maxAgendaPageSize
	}

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.list_time_off.tx_start_error", "err", err)
		return ListTimeOffOutput{}, utils.InternalError("")
	}
	defer func() {
		if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
			utils.SetSpanError(ctx, rollbackErr)
			logger.Error("photo_session.list_time_off.tx_rollback_error", "err", rollbackErr)
		}
	}()

	entries, err := s.repo.ListEntriesByRange(ctx, tx, input.PhotographerID, utils.ConvertToUTC(fromLocal), utils.ConvertToUTC(toLocal))
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.list_time_off.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListTimeOffOutput{}, utils.InternalError("")
	}

	timeOffEntries := make([]photosessionmodel.AgendaEntryInterface, 0)
	for _, entry := range entries {
		if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff {
			continue
		}
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
		timeOffEntries = append(timeOffEntries, entry)
	}

	sort.Slice(timeOffEntries, func(i, j int) bool {
		if timeOffEntries[i].StartsAt().Equal(timeOffEntries[j].StartsAt()) {
			return timeOffEntries[i].ID() < timeOffEntries[j].ID()
		}
		return timeOffEntries[i].StartsAt().Before(timeOffEntries[j].StartsAt())
	})

	total := len(timeOffEntries)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListTimeOffOutput{
		TimeOffs: timeOffEntries[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}
