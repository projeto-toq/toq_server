package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVisitInput holds the data to request a visit slot.
type CreateVisitInput struct {
	ListingIdentityID int64
	ScheduledStart    time.Time
	ScheduledEnd      time.Time
	Type              listingmodel.VisitMode
	RealtorNotes      string
	Source            string
}

func (s *visitService) CreateVisit(ctx context.Context, input CreateVisitInput) (visit listingmodel.VisitInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if !input.ScheduledStart.Before(input.ScheduledEnd) {
		return nil, utils.ValidationError("scheduledTime", "start must be before end")
	}

	requesterID, uidErr := s.globalService.GetUserIDFromContext(ctx)
	if uidErr != nil {
		return nil, uidErr
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("visit.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("visit.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	listingIdentity, err := s.listingRepo.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Listing")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.get_listing_identity_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	activeVersion, err := s.listingRepo.GetActiveListingVersion(ctx, tx, input.ListingIdentityID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.get_active_version_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.get_agenda_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	if err := s.validateWindow(ctx, tx, agenda, input); err != nil {
		return nil, err
	}

	// Basic conflict detection using blocking entries.
	entries, err := s.scheduleRepo.ListEntriesBetween(ctx, tx, agenda.ID(), input.ScheduledStart, input.ScheduledEnd)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.list_entries_error", "agenda_id", agenda.ID(), "err", err)
		return nil, utils.InternalError("")
	}
	for _, e := range entries {
		if !e.Blocking() {
			continue
		}
		if e.StartsAt().Before(input.ScheduledEnd) && e.EndsAt().After(input.ScheduledStart) {
			return nil, utils.ConflictError("Schedule conflict for requested interval")
		}
	}

	visit = listingmodel.NewVisit()
	visit.SetListingIdentityID(input.ListingIdentityID)
	visit.SetListingVersion(activeVersion.Version())
	visit.SetRequesterUserID(requesterID)
	visit.SetOwnerUserID(listingIdentity.UserID)
	visit.SetScheduledStart(input.ScheduledStart)
	visit.SetScheduledEnd(input.ScheduledEnd)
	visit.SetDurationMinutes(int64(input.ScheduledEnd.Sub(input.ScheduledStart).Minutes()))
	visit.SetStatus(listingmodel.VisitStatusPending)
	visit.SetType(input.Type)
	visit.SetCreatedBy(requesterID)
	if input.Source != "" {
		visit.SetSource(input.Source)
	}
	if input.RealtorNotes != "" {
		visit.SetRealtorNotes(input.RealtorNotes)
	}

	visitID, err := s.visitRepo.InsertVisit(ctx, tx, visit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.insert_visit_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}
	visit.SetID(visitID)

	if err = s.ensureVisitEntries(ctx, tx, agenda, visit, schedulemodel.EntryTypeVisitPending, true); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("visit.create.ensure_entries_error", "agenda_id", agenda.ID(), "visit_id", visitID, "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("visit.create.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}
	committed = true

	s.notifyVisitRequested(ctx, visit)

	return visit, nil
}
