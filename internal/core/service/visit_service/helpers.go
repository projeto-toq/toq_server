package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// loadVisitAndEntry retrieves the visit and its agenda entry (if any).
func (s *visitService) loadVisitAndEntry(ctx context.Context, tx *sql.Tx, visitID int64) (listingmodel.VisitInterface, schedulemodel.AgendaEntryInterface, bool, error) {
	visit, err := s.visitRepo.GetVisitByID(ctx, tx, visitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, false, utils.NotFoundError("Visit")
		}
		return nil, nil, false, err
	}

	entry, err := s.scheduleRepo.GetEntryByVisitID(ctx, tx, uint64(visitID))
	if err != nil {
		if err == sql.ErrNoRows {
			return visit, nil, false, nil
		}
		return nil, nil, false, err
	}

	return visit, entry, true, nil
}

func markFirstOwnerActionIfEmpty(visit listingmodel.VisitInterface, now time.Time) {
	if _, ok := visit.FirstOwnerActionAt(); !ok {
		visit.SetFirstOwnerActionAt(now)
	}
}
