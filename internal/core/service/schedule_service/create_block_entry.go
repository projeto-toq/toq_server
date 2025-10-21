package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) CreateBlockEntry(ctx context.Context, input CreateBlockEntryInput) (schedulemodel.AgendaEntryInterface, error) {
	if input.ListingID <= 0 {
		return nil, utils.ValidationError("listingId", "listingId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return nil, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if input.ActorID <= 0 {
		return nil, utils.ValidationError("actorId", "actorId must be greater than zero")
	}
	if !isBlockEntryType(input.EntryType) {
		return nil, utils.ValidationError("entryType", "entry type must be BLOCK or TEMP_BLOCK")
	}
	if !input.StartsAt.Before(input.EndsAt) {
		return nil, utils.ValidationError("range", "startsAt must be before endsAt")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.create_block_entry.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.create_block_entry.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, input.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_block_entry.get_agenda_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}

	if agenda.OwnerID() != input.OwnerID {
		return nil, utils.AuthorizationError("Owner does not match listing agenda")
	}

	existing, err := s.scheduleRepo.ListEntriesBetween(ctx, tx, agenda.ID(), input.StartsAt, input.EndsAt)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_block_entry.list_entries_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}
	for _, entry := range existing {
		if entry.Blocking() {
			return nil, utils.ConflictError("Agenda already blocked for this period")
		}
	}

	domain := schedulemodel.NewAgendaEntry()
	domain.SetAgendaID(agenda.ID())
	domain.SetEntryType(input.EntryType)
	domain.SetStartsAt(input.StartsAt)
	domain.SetEndsAt(input.EndsAt)
	domain.SetBlocking(true)
	if reason := strings.TrimSpace(input.Reason); reason != "" {
		domain.SetReason(reason)
	}

	id, err := s.scheduleRepo.InsertEntry(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_block_entry.insert_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}
	domain.SetID(id)

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.create_block_entry.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true

	return domain, nil
}
