package listingservices

import (
	"context"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type slotSortOption struct {
	column    string
	direction string
}

func (ls *listingService) ListPhotographerSlots(ctx context.Context, input ListPhotographerSlotsInput) (output ListPhotographerSlotsOutput, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return output, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

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

	sortOpt, sortErr := normalizeSlotSort(input.Sort)
	if sortErr != nil {
		return output, sortErr
	}

	params := photosessionmodel.SlotListParams{
		From:          input.From,
		To:            input.To,
		Period:        input.Period,
		Limit:         size,
		Offset:        (page - 1) * size,
		SortColumn:    sortOpt.column,
		SortDirection: sortOpt.direction,
	}

	tx, txErr := ls.gsi.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("listing.photo_session.list.start_tx_error", "err", txErr)
		return output, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := ls.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("listing.photo_session.list.rollback_error", "err", rbErr)
			}
		}
	}()

	slots, total, repoErr := ls.photoSessionRepo.ListAvailableSlots(ctx, tx, params)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("listing.photo_session.list.query_error", "err", repoErr)
		return output, utils.InternalError("")
	}

	if cmErr := ls.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("listing.photo_session.list.commit_error", "err", cmErr)
		return output, utils.InternalError("")
	}

	output = ListPhotographerSlotsOutput{
		Slots: slots,
		Total: total,
		Page:  page,
		Size:  size,
	}

	return output, nil
}

func normalizeSlotSort(sort string) (slotSortOption, error) {
	switch sort {
	case "", "date_asc":
		return slotSortOption{column: "slot_date", direction: "ASC"}, nil
	case "date_desc":
		return slotSortOption{column: "slot_date", direction: "DESC"}, nil
	case "photographer_asc":
		return slotSortOption{column: "photographer_user_id", direction: "ASC"}, nil
	case "photographer_desc":
		return slotSortOption{column: "photographer_user_id", direction: "DESC"}, nil
	default:
		return slotSortOption{}, utils.BadRequest("unsupported sort parameter")
	}
}

func sessionWindowForPeriod(date time.Time, period photosessionmodel.SlotPeriod) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	if period == photosessionmodel.SlotPeriodAfternoon {
		start = time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, time.UTC)
	}
	end := start.Add(4 * time.Hour)
	return start, end
}

func listingEligibleForPhotoSession(status listingmodel.ListingStatus) bool {
	switch status {
	case listingmodel.StatusDraft, listingmodel.StatusAwaitingPhoto:
		return true
	default:
		return false
	}
}
