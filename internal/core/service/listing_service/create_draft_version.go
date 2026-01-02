package listingservices

import (
	"context"
	"database/sql"
	"errors"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateDraftVersionInput carries the data required to create a new draft version.
type CreateDraftVersionInput struct {
	ListingIdentityID int64
}

// CreateDraftVersionOutput contains the created draft version information.
type CreateDraftVersionOutput struct {
	VersionID int64
	Version   uint8
	Status    string
}

// CreateDraftVersion creates a new draft version from an existing active version.
// It validates the active version's status and either returns an existing draft or creates a new one.
func (ls *listingService) CreateDraftVersion(ctx context.Context, input CreateDraftVersionInput) (output CreateDraftVersionOutput, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := ls.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.create_draft.tx_start_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.create_draft.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	output, err = ls.createDraftVersion(ctx, tx, input)
	if err != nil {
		return
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.create_draft.tx_commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}

	logger.Info("listing.create_draft.success", "listing_identity_id", input.ListingIdentityID, "version_id", output.VersionID, "version", output.Version)
	return
}

func (ls *listingService) createDraftVersion(ctx context.Context, tx *sql.Tx, input CreateDraftVersionInput) (output CreateDraftVersionOutput, err error) {
	logger := utils.LoggerFromContext(ctx)

	// Check if draft version already exists (idempotency)
	draftVersion, draftErr := ls.listingRepository.GetDraftVersionByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if draftErr != nil && !errors.Is(draftErr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, draftErr)
		logger.Error("listing.create_draft.get_draft_error", "err", draftErr, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	if draftVersion != nil {
		// Return existing draft
		output.VersionID = draftVersion.ID()
		output.Version = draftVersion.Version()
		output.Status = draftVersion.Status().String()
		logger.Info("listing.create_draft.existing_draft", "listing_identity_id", input.ListingIdentityID, "version_id", output.VersionID)
		return output, nil
	}

	// Get active version
	activeVersion, activeErr := ls.listingRepository.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if activeErr != nil {
		if errors.Is(activeErr, sql.ErrNoRows) {
			return output, utils.NotFoundError("Active listing version not found")
		}
		utils.SetSpanError(ctx, activeErr)
		logger.Error("listing.create_draft.get_active_error", "err", activeErr, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	// Validate ownership
	userID, uidErr := ls.gsi.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return output, uidErr
	}

	if activeVersion.UserID() != userID {
		logger.Warn("unauthorized_create_draft_attempt",
			"listing_identity_id", input.ListingIdentityID,
			"listing_version_id", activeVersion.ID(),
			"requester_user_id", userID,
			"owner_user_id", activeVersion.UserID())
		return output, utils.AuthorizationError("not the listing owner")
	}

	// Validate active version status
	if validErr := ls.validateStatusForDraftCreation(activeVersion.Status()); validErr != nil {
		return output, validErr
	}

	// Create new draft version
	newVersion := listingmodel.NewListingVersion()
	newVersion.SetListingIdentityID(input.ListingIdentityID)
	newVersion.SetListingUUID(activeVersion.ListingUUID())
	newVersion.SetUserID(activeVersion.UserID())
	newVersion.SetCode(activeVersion.Code())
	newVersion.SetVersion(activeVersion.Version() + 1)
	newVersion.SetStatus(listingmodel.StatusDraft)
	newVersion.SetDeleted(false)

	// Copy all fields from active version
	ls.copyVersionFields(activeVersion, newVersion)

	// Create version record
	if createErr := ls.listingRepository.CreateListingVersion(ctx, tx, newVersion); createErr != nil {
		utils.SetSpanError(ctx, createErr)
		logger.Error("listing.create_draft.create_version_error", "err", createErr, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	// Clone satellite entities
	if cloneErr := ls.listingRepository.CloneListingVersionSatellites(ctx, tx, activeVersion.ID(), newVersion.ID()); cloneErr != nil {
		utils.SetSpanError(ctx, cloneErr)
		logger.Error("listing.create_draft.clone_satellites_error", "err", cloneErr, "source", activeVersion.ID(), "target", newVersion.ID())
		return output, utils.InternalError("")
	}

	output.VersionID = newVersion.ID()
	output.Version = newVersion.Version()
	output.Status = newVersion.Status().String()

	return output, nil
}

func (ls *listingService) validateStatusForDraftCreation(status listingmodel.ListingStatus) *utils.HTTPError {
	if listingmodel.StatusAllowsDraftClone(status) {
		return nil
	}

	switch status {
	case listingmodel.StatusPublished:
		return utils.ConflictError("Listing is published. Suspend it via status update before creating a draft version")
	case listingmodel.StatusPendingAdminReview, listingmodel.StatusPendingOwnerApproval:
		return utils.NewHTTPErrorWithSource(423, "Listing is locked in workflow and cannot be copied")
	case listingmodel.StatusExpired, listingmodel.StatusArchived, listingmodel.StatusClosed:
		return utils.NewHTTPErrorWithSource(410, "Listing is permanently closed and cannot be edited")
	default:
		return utils.BadRequest("Listing status does not allow draft creation")
	}
}

func (ls *listingService) copyVersionFields(source, target listingmodel.ListingVersionInterface) {
	if source.HasTitle() {
		target.SetTitle(source.Title())
	}
	// Copy required address fields
	target.SetStreet(source.Street())
	target.SetNumber(source.Number())
	target.SetComplement(source.Complement())
	target.SetNeighborhood(source.Neighborhood())
	target.SetCity(source.City())
	target.SetState(source.State())
	target.SetZipCode(source.ZipCode())
	if source.HasComplex() {
		target.SetComplex(source.Complex())
	} else {
		target.UnsetComplex()
	}
	target.SetListingType(source.ListingType())

	if source.HasOwner() {
		target.SetOwner(source.Owner())
	}
	if source.HasLandSize() {
		target.SetLandSize(source.LandSize())
	}
	if source.HasCorner() {
		target.SetCorner(source.Corner())
	}
	if source.HasNonBuildable() {
		target.SetNonBuildable(source.NonBuildable())
	}
	if source.HasBuildable() {
		target.SetBuildable(source.Buildable())
	}
	if source.HasDelivered() {
		target.SetDelivered(source.Delivered())
	}
	if source.HasWhoLives() {
		target.SetWhoLives(source.WhoLives())
	}
	if source.HasDescription() {
		target.SetDescription(source.Description())
	}
	if source.HasTransaction() {
		target.SetTransaction(source.Transaction())
	}
	if source.HasSellNet() {
		target.SetSellNet(source.SellNet())
	}
	if source.HasRentNet() {
		target.SetRentNet(source.RentNet())
	}
	if source.HasCondominium() {
		target.SetCondominium(source.Condominium())
	}
	if source.HasAnnualTax() {
		target.SetAnnualTax(source.AnnualTax())
	}
	if source.HasMonthlyTax() {
		target.SetMonthlyTax(source.MonthlyTax())
	}
	if source.HasAnnualGroundRent() {
		target.SetAnnualGroundRent(source.AnnualGroundRent())
	}
	if source.HasMonthlyGroundRent() {
		target.SetMonthlyGroundRent(source.MonthlyGroundRent())
	}
	if source.HasExchange() {
		target.SetExchange(source.Exchange())
	}
	if source.HasExchangePercentual() {
		target.SetExchangePercentual(source.ExchangePercentual())
	}
	if source.HasInstallment() {
		target.SetInstallment(source.Installment())
	}
	if source.HasFinancing() {
		target.SetFinancing(source.Financing())
	}
	if source.HasVisit() {
		target.SetVisit(source.Visit())
	}
	if source.HasTenantName() {
		target.SetTenantName(source.TenantName())
	}
	if source.HasTenantEmail() {
		target.SetTenantEmail(source.TenantEmail())
	}
	if source.HasTenantPhone() {
		target.SetTenantPhone(source.TenantPhone())
	}
	if source.HasAccompanying() {
		target.SetAccompanying(source.Accompanying())
	}
}
