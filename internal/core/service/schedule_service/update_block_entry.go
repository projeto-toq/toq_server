package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) UpdateBlockEntry(ctx context.Context, input UpdateBlockEntryInput) (schedulemodel.AgendaEntryInterface, error) {
	if input.EntryID == 0 {
		return nil, utils.ValidationError("entryId", "entryId must be greater than zero")
	}
	if input.ListingIdentityID <= 0 {
		return nil, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
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
	reqLoc, tzErr := utils.ResolveLocation("timezone", input.Timezone)
	if tzErr != nil {
		return nil, tzErr
	}

	startsAtLocal := input.StartsAt.In(reqLoc)
	endsAtLocal := input.EndsAt.In(reqLoc)
	if !startsAtLocal.Before(endsAtLocal) {
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
		logger.Error("schedule.update_block_entry.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.update_block_entry.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.update_block_entry.get_agenda_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	if agenda.OwnerID() != input.OwnerID {
		return nil, utils.AuthorizationError("Owner does not match listing agenda")
	}

	agendaLoc, agendaTzErr := utils.ResolveLocation("timezone", agenda.Timezone())
	if agendaTzErr != nil {
		return nil, agendaTzErr
	}

	startsAtAgenda := startsAtLocal.In(agendaLoc)
	endsAtAgenda := endsAtLocal.In(agendaLoc)
	startsAtUTC := startsAtAgenda.UTC()
	endsAtUTC := endsAtAgenda.UTC()

	existing, err := s.scheduleRepo.GetEntryByID(ctx, tx, input.EntryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda entry")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.update_block_entry.get_entry_error", "entry_id", input.EntryID, "err", err)
		return nil, utils.InternalError("")
	}

	if existing.AgendaID() != agenda.ID() {
		return nil, utils.AuthorizationError("Entry does not belong to provided listing")
	}

	if !isBlockEntryType(existing.EntryType()) {
		return nil, utils.ConflictError("Only blocking entries can be updated through this operation")
	}

	overlaps, err := s.scheduleRepo.ListEntriesBetween(ctx, tx, agenda.ID(), startsAtUTC, endsAtUTC)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.update_block_entry.list_entries_error", "entry_id", input.EntryID, "err", err)
		return nil, utils.InternalError("")
	}
	for _, entry := range overlaps {
		if entry.ID() == existing.ID() {
			continue
		}
		if entry.Blocking() {
			return nil, utils.ConflictError("Agenda already blocked for this period")
		}
	}

	existing.SetEntryType(input.EntryType)
	existing.SetStartsAt(startsAtUTC)
	existing.SetEndsAt(endsAtUTC)
	existing.SetBlocking(true)
	if reason := strings.TrimSpace(input.Reason); reason != "" {
		existing.SetReason(reason)
	} else {
		existing.ClearReason()
	}

	if err := s.scheduleRepo.UpdateEntry(ctx, tx, existing); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.update_block_entry.update_error", "entry_id", input.EntryID, "err", err)
		return nil, utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.update_block_entry.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true
	existing.SetStartsAt(startsAtLocal)
	existing.SetEndsAt(endsAtLocal)

	return existing, nil
}
