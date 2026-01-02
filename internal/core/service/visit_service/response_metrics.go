package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// recordOwnerResponseMetrics sets the first owner action timestamp and updates listing-level metrics once.
func (s *visitService) recordOwnerResponseMetrics(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface, actorID int64, actionTime time.Time) error {
	if actorID != visit.OwnerUserID() {
		return nil
	}

	if _, ok := visit.FirstOwnerActionAt(); ok {
		return nil
	}

	visit.SetFirstOwnerActionAt(actionTime)

	requestedAt := visit.RequestedAt()
	if requestedAt.IsZero() {
		requestedAt = visit.ScheduledStart()
	}

	delta := actionTime.Sub(requestedAt)
	if delta < 0 {
		delta = 0
	}
	maxDelta := 24 * time.Hour * 365
	if delta > maxDelta {
		delta = maxDelta
	}

	deltaSeconds := int64(delta / time.Second)
	if err := s.listingRepo.UpdateOwnerResponseStats(ctx, tx, visit.ListingIdentityID(), deltaSeconds, actionTime); err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("ListingIdentity")
		}
		return err
	}

	return nil
}
