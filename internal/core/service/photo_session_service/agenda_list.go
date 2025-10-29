package photosessionservices

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListAgenda returns agenda entries within the interval.
func (s *photoSessionService) ListAgenda(ctx context.Context, input ListAgendaInput) (ListAgendaOutput, error) {
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.ListAgenda")
	if err != nil {
		return ListAgendaOutput{}, derrors.Infra("failed to generate tracer", err)
	}
	defer spanEnd()

	if err := validateListAgendaInput(input); err != nil {
		return ListAgendaOutput{}, err
	}

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return ListAgendaOutput{}, tzErr
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

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.tx_start_error", "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to start transaction")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			utils.LoggerFromContext(ctx).Error("service.list_agenda.tx_rollback_error", "err", rbErr)
		}
	}()

	entries, err := s.repo.ListEntriesByRange(ctx, tx, input.PhotographerID, input.StartDate.UTC(), input.EndDate.UTC())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("service.list_agenda.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return ListAgendaOutput{}, derrors.Wrap(err, derrors.KindInfra, "failed to list agenda entries")
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].StartsAt().Equal(entries[j].StartsAt()) {
			return entries[i].ID() < entries[j].ID()
		}
		return entries[i].StartsAt().Before(entries[j].StartsAt())
	})

	slots := make([]AgendaSlot, 0, len(entries))
	for _, entry := range entries {
		slots = append(slots, s.buildAgendaSlot(entry, loc))
	}

	total := len(slots)
	start := (page - 1) * size
	if start > total {
		start = total
	}
	end := start + size
	if end > total {
		end = total
	}

	return ListAgendaOutput{
		Slots:    slots[start:end],
		Total:    int64(total),
		Page:     page,
		Size:     size,
		Timezone: loc.String(),
	}, nil
}

func (s *photoSessionService) buildAgendaSlot(entry photosessionmodel.AgendaEntryInterface, loc *time.Location) AgendaSlot {
	start := entry.StartsAt().In(loc)
	end := entry.EndsAt().In(loc)

	slot := AgendaSlot{
		EntryID:        entry.ID(),
		PhotographerID: entry.PhotographerUserID(),
		Start:          start,
		End:            end,
		Status:         photosessionmodel.SlotStatusBlocked,
		GroupID:        buildAgendaGroupID(entry, loc),
		Source:         entry.Source(),
		IsHoliday:      entry.EntryType() == photosessionmodel.AgendaEntryTypeHoliday,
		IsTimeOff:      entry.EntryType() == photosessionmodel.AgendaEntryTypeTimeOff,
		Timezone:       loc.String(),
		EntryType:      entry.EntryType(),
	}

	if !entry.Blocking() {
		slot.Status = photosessionmodel.SlotStatusAvailable
	}

	switch entry.EntryType() {
	case photosessionmodel.AgendaEntryTypePhotoSession:
		slot.Status = photosessionmodel.SlotStatusBooked
		if sourceID, ok := entry.SourceID(); ok && sourceID != nil {
			slot.SourceID = *sourceID
		}
	case photosessionmodel.AgendaEntryTypeHoliday:
		if sourceID, ok := entry.SourceID(); ok && sourceID != nil {
			slot.HolidayCalendarIDs = []uint64{*sourceID}
		}
		if reason, ok := entry.Reason(); ok {
			slot.HolidayLabels = []string{reason}
		}
	case photosessionmodel.AgendaEntryTypeTimeOff:
		if reason, ok := entry.Reason(); ok {
			slot.Reason = reason
		}
	}

	return slot
}

func buildAgendaGroupID(entry photosessionmodel.AgendaEntryInterface, loc *time.Location) string {
	dayKey := entry.StartsAt().In(loc).Format("2006-01-02")
	return fmt.Sprintf("%s-%s", strings.ToLower(string(entry.Source())), dayKey)
}
