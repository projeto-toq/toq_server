package listingservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ls *listingService) ListPhotographerSlots(ctx context.Context, input ListPhotographerSlotsInput) (output ListPhotographerSlotsOutput, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, userErr := ls.gsi.GetUserIDFromContext(ctx)
	if userErr != nil {
		return output, userErr
	}

	if input.ListingIdentityID <= 0 {
		return output, utils.BadRequest("listingIdentityId must be greater than zero")
	}

	loc := input.Location
	if loc == nil {
		return output, utils.BadRequest("timezone is required")
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.photo_session.list_slots.start_ro_tx_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if tx != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.list_slots.ro_rollback_error", "err", rbErr)
			}
		}
	}()

	// Load active listing version to validate eligibility
	listing, repoErr := ls.listingRepository.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if repoErr != nil {
		if errors.Is(repoErr, sql.ErrNoRows) {
			return output, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.photo_session.list_slots.get_listing_error", "err", repoErr, "listing_identity_id", input.ListingIdentityID)
		return output, utils.InternalError("")
	}

	if listing.Deleted() {
		return output, utils.BadRequest("listing is not available")
	}

	if listing.UserID() != userID {
		return output, utils.AuthorizationError("listing does not belong to current user")
	}

	if !listingEligibleForPhotoSession(listing.Status()) {
		return output, derrors.ErrListingNotEligible
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.list_slots.ro_commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}
	tx = nil

	page := input.Page
	if page <= 0 {
		page = defaultSlotsPage
	}

	size := input.Size
	if size <= 0 {
		size = defaultSlotsPageSize
	}
	if size > maxSlotsPageSize {
		size = maxSlotsPageSize
	}

	if input.From != nil && input.To != nil && input.From.After(*input.To) {
		return output, utils.BadRequest("from must be before to")
	}

	sortKey, sortErr := normalizeSlotSort(strings.TrimSpace(input.Sort))
	if sortErr != nil {
		return output, sortErr
	}

	availabilityInput := photosessionservices.ListAvailabilityInput{
		From:              input.From,
		To:                input.To,
		Period:            input.Period,
		Page:              page,
		Size:              size,
		Sort:              sortKey,
		Location:          loc,
		ListingIdentityID: listing.IdentityID(),
	}

	availability, svcErr := ls.photoSessionSvc.ListAvailability(ctx, availabilityInput)
	if svcErr != nil {
		return output, svcErr
	}

	slots := make([]photosessionmodel.PhotographerSlotInterface, 0, len(availability.Slots))
	for _, slot := range availability.Slots {
		slots = append(slots, availabilitySlotToPhotographerSlot(slot))
	}

	output = ListPhotographerSlotsOutput{
		Slots: slots,
		Total: availability.Total,
		Page:  page,
		Size:  size,
	}

	return output, nil
}

func normalizeSlotSort(sort string) (string, error) {
	switch sort {
	case "", "start_asc", "date_asc":
		return "start_asc", nil
	case "start_desc", "date_desc":
		return "start_desc", nil
	case "photographer_asc":
		return "photographer_asc", nil
	case "photographer_desc":
		return "photographer_desc", nil
	default:
		return "", utils.BadRequest("unsupported sort parameter")
	}
}

func listingEligibleForPhotoSession(status listingmodel.ListingStatus) bool {
	switch status {
	case
		// listingmodel.StatusDraft,
		listingmodel.StatusPendingPhotoScheduling,
		listingmodel.StatusPendingPhotoConfirmation,
		listingmodel.StatusPhotosScheduled: //,
		// listingmodel.StatusPendingPhotoProcessing,
		// listingmodel.StatusPendingOwnerApproval		:
		return true
	default:
		return false
	}
}

func availabilitySlotToPhotographerSlot(slot photosessionservices.AvailabilitySlot) photosessionmodel.PhotographerSlotInterface {
	ps := photosessionmodel.NewPhotographerSlot()
	ps.SetID(slot.SlotID)
	ps.SetPhotographerUserID(slot.PhotographerID)
	ps.SetSlotStart(slot.Start)
	ps.SetSlotEnd(slot.End)
	ps.SetSlotDate(time.Date(slot.Start.Year(), slot.Start.Month(), slot.Start.Day(), 0, 0, 0, 0, slot.Start.Location()))
	ps.SetStatus(photosessionmodel.SlotStatusAvailable)
	ps.SetPeriod(slot.Period)
	return ps
}
