package listingservices

import (
	"context"
	"database/sql"
	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) PromoteListingVersion(ctx context.Context, input PromoteListingVersionInput) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.VersionID <= 0 {
		return utils.ValidationError("versionId", "versionId must be greater than zero")
	}

	tx, err := ls.gsi.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.tx_start_error", "err", err, "version_id", input.VersionID)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.promote.tx_rollback_error", "err", rbErr, "version_id", input.VersionID)
			}
		}
	}()

	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return uidErr
	}

	listingVersionID := input.VersionID

	snapshot, err := ls.listingRepository.GetListingForEndUpdate(ctx, tx, listingVersionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.fetch_error", "err", err, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	if snapshot.UserID != userID {
		return utils.AuthorizationError("Only listing owner can promote version")
	}

	if snapshot.Status != listingmodel.StatusDraft {
		return utils.ConflictError("Only draft versions can be promoted")
	}

	if verr := ls.validateListingBeforeEndUpdate(ctx, tx, snapshot); verr != nil {
		return verr
	}

	// Retrieve all versions for this identity to find the current active version
	versionSummaries, listErr := ls.listingRepository.ListListingVersions(ctx, tx, listingrepository.ListListingVersionsFilter{
		ListingIdentityID: snapshot.ListingID,
		IncludeDeleted:    false,
	})
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("listing.promote.list_versions_error", "err", listErr, "listing_identity_id", snapshot.ListingID)
		return utils.InternalError("")
	}

	var currentActiveVersion listingmodel.ListingVersionInterface
	for _, summary := range versionSummaries {
		if summary.IsActive && summary.Version != nil {
			currentActiveVersion = summary.Version
			break
		}
	}

	// Determine target status based on current active version status
	targetStatus := listingmodel.StatusPendingAvailability
	if currentActiveVersion != nil {
		targetStatus = currentActiveVersion.Status()
	}

	updateErr := ls.listingRepository.UpdateListingStatus(ctx, tx, listingVersionID, targetStatus, listingmodel.StatusDraft)
	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return utils.ConflictError("Listing status changed while promoting version")
		}
		utils.SetSpanError(ctx, updateErr)
		logger.Error("listing.promote.update_status_error", "err", updateErr, "listing_version_id", listingVersionID, "target_status", targetStatus)
		return utils.InternalError("")
	}

	// Set the new active version
	if setErr := ls.listingRepository.SetListingActiveVersion(ctx, tx, snapshot.ListingID, listingVersionID); setErr != nil {
		utils.SetSpanError(ctx, setErr)
		logger.Error("listing.promote.set_active_error", "err", setErr, "listing_identity_id", snapshot.ListingID, "version_id", listingVersionID)
		return utils.InternalError("")
	}

	if auditErr := ls.gsi.CreateAudit(ctx, tx, globalmodel.TableListings, "Versão de anúncio promovida"); auditErr != nil {
		return auditErr
	}

	// If target status is PendingAvailability (meaning this is the first version or we're starting fresh), create agenda
	if targetStatus == listingmodel.StatusPendingAvailability {
		timezone := resolveListingTimezone(snapshot)
		agendaInput := scheduleservices.CreateDefaultAgendaInput{
			ListingIdentityID: snapshot.ListingID,
			OwnerID:           userID,
			Timezone:          timezone,
			ActorID:           userID,
		}
		if _, agendaErr := ls.scheduleService.CreateDefaultAgendaWithTx(ctx, tx, agendaInput); agendaErr != nil {
			utils.SetSpanError(ctx, agendaErr)
			logger.Error("listing.promote.create_default_agenda_error", "err", agendaErr, "listing_identity_id", snapshot.ListingID)
			return agendaErr
		}
	}

	if err = ls.gsi.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("listing.promote.tx_commit_error", "err", err, "listing_version_id", listingVersionID)
		return utils.InternalError("")
	}

	logger.Info("listing.promote.completed", "listing_version_id", listingVersionID, "listing_identity_id", snapshot.ListingID, "new_status", targetStatus.String())

	return nil
}
