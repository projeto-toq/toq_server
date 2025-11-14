package listingservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) DiscardDraftVersion(ctx context.Context, input DiscardDraftVersionInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.tx_start_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.discard.tx_rollback_error", "err", rbErr, "version_id", input.VersionID)
			}
		}
	}()

	draft, err := ls.listingRepository.GetListingVersionByID(ctx, tx, input.VersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing version")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.get_version_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}

	if draft.Status() != listingmodel.StatusDraft {
		return utils.ConflictError("Only draft versions can be discarded")
	}

	if draft.ActiveVersionID() == draft.ID() {
		return utils.ConflictError("Cannot discard the active listing version")
	}

	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	if draft.UserID() != userID {
		return utils.AuthorizationError("Only listing owner can discard a draft version")
	}

	draft.SetDeleted(true)

	if updateErr := ls.listingRepository.UpdateListingVersion(ctx, tx, draft); updateErr != nil && !errors.Is(updateErr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.discard.update_error", "err", updateErr, "version_id", draft.ID())
		return utils.InternalError("")
	}

	if auditErr := ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Versão de anúncio descartada"); auditErr != nil {
		return auditErr
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.discard.tx_commit_error", "err", err, "version_id", draft.ID())
		return utils.InternalError("")
	}

	logger.Info("listing.discard.completed", "version_id", draft.ID(), "identity_id", draft.IdentityID())

	return nil
}
