package listingservices

import (
	"context"
	"time"

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

	sortKey, sortErr := normalizeSlotSort(input.Sort)
	if sortErr != nil {
		return output, sortErr
	}

	availabilityInput := photosessionservices.ListAvailabilityInput{
		From:   input.From,
		To:     input.To,
		Period: input.Period,
		Page:   page,
		Size:   size,
		Sort:   sortKey,
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
		listingmodel.StatusPendingAvailabilityConfirm,
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
