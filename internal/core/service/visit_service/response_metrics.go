package visitservice

import (
	"context"
	"database/sql"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	ownermetricsrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/owner_metrics_repository"
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

	ownerID := visit.OwnerUserID()
	if ownerID <= 0 {
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
	input := ownermetricsrepository.VisitResponseInput{
		OwnerID:      ownerID,
		DeltaSeconds: deltaSeconds,
		RespondedAt:  actionTime,
	}
	if err := s.ownerMetrics.UpsertVisitResponse(ctx, tx, input); err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFoundError("Owner")
		}
		return err
	}

	return nil
}
