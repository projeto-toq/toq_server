package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) ListBlockEntries(ctx context.Context, filter schedulemodel.BlockEntriesFilter) (schedulemodel.BlockEntriesResult, error) {
	if filter.OwnerID <= 0 {
		return schedulemodel.BlockEntriesResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if filter.ListingID <= 0 {
		return schedulemodel.BlockEntriesResult{}, utils.ValidationError("listingId", "listingId must be greater than zero")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return schedulemodel.BlockEntriesResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.BlockEntriesResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.list_block_entries.tx_start_error", "err", txErr, "listing_id", filter.ListingID)
		return schedulemodel.BlockEntriesResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.list_block_entries.tx_rollback_error", "err", rbErr, "listing_id", filter.ListingID)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, filter.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedulemodel.BlockEntriesResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_block_entries.get_agenda_error", "err", err, "listing_id", filter.ListingID)
		return schedulemodel.BlockEntriesResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != filter.OwnerID {
		return schedulemodel.BlockEntriesResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	agendaLoc, tzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if tzErr != nil {
		return schedulemodel.BlockEntriesResult{}, tzErr
	}

	loc := filter.Range.Loc
	if loc == nil {
		loc = agendaLoc
	}
	if loc == nil {
		loc = time.UTC
	}

	repoFilter := filter
	repoFilter.Range.From, repoFilter.Range.To = utils.NormalizeRangeToUTC(filter.Range.From, filter.Range.To, loc)
	repoFilter.Range.Loc = time.UTC

	pageResult, err := s.scheduleRepo.ListBlockEntries(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_block_entries.repo_error", "err", err, "listing_id", filter.ListingID)
		return schedulemodel.BlockEntriesResult{}, utils.InternalError("")
	}

	items := make([]schedulemodel.AgendaEntryInterface, 0, len(pageResult.Entries))
	for _, entry := range pageResult.Entries {
		if entry == nil {
			continue
		}
		entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
		entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))
		items = append(items, entry)
	}

	return schedulemodel.BlockEntriesResult{
		Items:    items,
		Total:    pageResult.Total,
		Timezone: loc.String(),
	}, nil
}
