package visitservice

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	ownerBlockReason   = "VISIT_OWNER_BLOCK"
	realtorBlockReason = "VISIT_REALTOR_BLOCK"
)

// loadVisit resolves a visit by ID with a not-found mapping.
func (s *visitService) loadVisit(ctx context.Context, tx *sql.Tx, visitID int64) (listingmodel.VisitInterface, error) {
	visit, err := s.visitRepo.GetVisitByID(ctx, tx, visitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("Visit")
		}
		return nil, err
	}
	return visit, nil
}

// listVisitEntries fetches agenda entries for the visit window and filters by VisitID.
func (s *visitService) listVisitEntries(ctx context.Context, tx *sql.Tx, agendaID uint64, visit listingmodel.VisitInterface) ([]schedulemodel.AgendaEntryInterface, error) {
	entries, err := s.scheduleRepo.ListEntriesBetween(ctx, tx, agendaID, visit.ScheduledStart(), visit.ScheduledEnd())
	if err != nil {
		return nil, err
	}

	filtered := make([]schedulemodel.AgendaEntryInterface, 0, len(entries))
	for _, e := range entries {
		if visitID, ok := e.VisitID(); ok && visitID == uint64(visit.ID()) {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// ensureVisitEntries creates or updates owner/realtor entries with the given type and blocking flag.
func (s *visitService) ensureVisitEntries(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface, visit listingmodel.VisitInterface, entryType schedulemodel.EntryType, blocking bool) error {
	existing, err := s.listVisitEntries(ctx, tx, agenda.ID(), visit)
	if err != nil {
		return err
	}

	byReason := map[string]schedulemodel.AgendaEntryInterface{}
	for _, e := range existing {
		if reason, ok := e.Reason(); ok {
			byReason[reason] = e
			continue
		}
		// Backfill legacy entries without reason as owner blocks.
		e.SetReason(ownerBlockReason)
		if err := s.scheduleRepo.UpdateEntry(ctx, tx, e); err != nil {
			return err
		}
		byReason[ownerBlockReason] = e
	}

	upsert := func(reason string) error {
		if entry, ok := byReason[reason]; ok {
			entry.SetEntryType(entryType)
			entry.SetBlocking(blocking)
			entry.SetStartsAt(visit.ScheduledStart())
			entry.SetEndsAt(visit.ScheduledEnd())
			return s.scheduleRepo.UpdateEntry(ctx, tx, entry)
		}

		entry := schedulemodel.NewAgendaEntry()
		entry.SetAgendaID(agenda.ID())
		entry.SetVisitID(uint64(visit.ID()))
		entry.SetEntryType(entryType)
		entry.SetBlocking(blocking)
		entry.SetStartsAt(visit.ScheduledStart())
		entry.SetEndsAt(visit.ScheduledEnd())
		entry.SetReason(reason)
		_, err := s.scheduleRepo.InsertEntry(ctx, tx, entry)
		return err
	}

	if err := upsert(ownerBlockReason); err != nil {
		return err
	}
	if err := upsert(realtorBlockReason); err != nil {
		return err
	}
	return nil
}

// removeVisitEntries deletes both owner and realtor entries for the visit if present.
func (s *visitService) removeVisitEntries(ctx context.Context, tx *sql.Tx, agenda schedulemodel.AgendaInterface, visit listingmodel.VisitInterface) error {
	entries, err := s.listVisitEntries(ctx, tx, agenda.ID(), visit)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if err := s.scheduleRepo.DeleteEntry(ctx, tx, e.ID()); err != nil {
			return err
		}
	}
	return nil
}
