package scheduleservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) DeleteBlockEntry(ctx context.Context, input DeleteEntryInput) error {
	if input.EntryID == 0 {
		return utils.ValidationError("entryId", "entryId must be greater than zero")
	}
	if input.ListingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.delete_block_entry.tx_start_error", "err", txErr)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.delete_block_entry.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.delete_block_entry.get_agenda_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return utils.InternalError("")
	}

	entry, err := s.scheduleRepo.GetEntryByID(ctx, tx, input.EntryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Agenda entry")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.delete_block_entry.get_entry_error", "entry_id", input.EntryID, "err", err)
		return utils.InternalError("")
	}

	if entry.AgendaID() != agenda.ID() {
		return utils.AuthorizationError("Entry does not belong to provided listing")
	}

	if !entry.Blocking() || !isBlockEntryType(entry.EntryType()) {
		return utils.ConflictError("Only blocking entries can be deleted through this operation")
	}

	if err := s.scheduleRepo.DeleteEntry(ctx, tx, input.EntryID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.delete_block_entry.delete_error", "entry_id", input.EntryID, "err", err)
		return utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.delete_block_entry.tx_commit_error", "err", cmErr)
		return utils.InternalError("")
	}

	committed = true
	return nil
}
