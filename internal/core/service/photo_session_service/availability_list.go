package photosessionservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListAvailability computes booking availability windows for photographers.
func (s *photoSessionService) ListAvailability(ctx context.Context, input ListAvailabilityInput) (ListAvailabilityOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAvailability")
	if err != nil {
		return ListAvailabilityOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return ListAvailabilityOutput{}, tzErr
	}

	page := input.Page
	if page <= 0 {
		page = defaultAgendaPage
	}

	size := input.Size
	if size <= 0 {
		size = defaultAgendaSize
	}
	if size > maxAgendaPageSize {
		size = maxAgendaPageSize
	}

	workdayStart := input.WorkdayStartHour
	if workdayStart <= 0 {
		workdayStart = defaultWorkdayStartHour
	}
	workdayEnd := input.WorkdayEndHour
	if workdayEnd <= 0 {
		workdayEnd = defaultWorkdayEndHour
	}
	if workdayEnd <= workdayStart {
		return ListAvailabilityOutput{}, derrors.Validation("workdayEndHour must be greater than workdayStartHour", nil)
	}

	now := s.now().In(loc)
	var rangeStart time.Time
	if input.From != nil {
		rangeStart = input.From.In(loc)
	} else {
		rangeStart = now
	}
	var rangeEnd time.Time
	if input.To != nil {
		rangeEnd = input.To.In(loc)
	} else {
		rangeEnd = rangeStart.AddDate(0, defaultHorizonMonths, 0)
	}
	if rangeEnd.Before(rangeStart) {
		return ListAvailabilityOutput{}, derrors.Validation("to must be after from", nil)
	}

	slotDuration := time.Duration(s.cfg.SlotDurationMinutes) * time.Minute
	if slotDuration <= 0 {
		slotDuration = 4 * time.Hour
	}

	filterPeriod := input.Period

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("photo_session.list_availability.tx_start_error", "err", txErr)
		return ListAvailabilityOutput{}, derrors.Infra("failed to start transaction", txErr)
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("photo_session.list_availability.tx_rollback_error", "err", rbErr)
		}
	}()

	photographerIDs, repoErr := s.repo.ListPhotographerIDs(ctx, tx)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("photo_session.list_availability.list_photographers_error", "err", repoErr)
		return ListAvailabilityOutput{}, derrors.Infra("failed to list photographers", repoErr)
	}

	availability := make([]AvailabilitySlot, 0)
	for _, photographerID := range photographerIDs {
		entries, err := s.repo.ListEntriesByRange(ctx, tx, photographerID, rangeStart.UTC(), rangeEnd.UTC())
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("photo_session.list_availability.list_entries_error", "photographer_id", photographerID, "err", err)
			return ListAvailabilityOutput{}, derrors.Infra("failed to load agenda entries", err)
		}

		workingRanges := buildWorkingRanges(rangeStart, rangeEnd, loc, workdayStart, workdayEnd)
		freeRanges := applyBlockingEntries(workingRanges, entries, loc)
		freeRanges = prunePastRanges(freeRanges, now)

		for _, free := range freeRanges {
			slots := splitIntoSlots(free, slotDuration)
			for _, slot := range slots {
				period := determineSlotPeriod(slot.start)
				if filterPeriod != nil && period != *filterPeriod {
					continue
				}
				id := encodeSlotID(photographerID, slot.start.UTC())
				availability = append(availability, AvailabilitySlot{
					SlotID:         id,
					PhotographerID: photographerID,
					Start:          slot.start,
					End:            slot.end,
					Period:         period,
					SourceTimezone: loc.String(),
				})
			}
		}
	}

	sortAvailabilitySlots(availability, input.Sort)

	total := len(availability)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListAvailabilityOutput{
		Slots:    availability[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}
