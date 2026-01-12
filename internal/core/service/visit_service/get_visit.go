package visitservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) GetVisit(ctx context.Context, visitID int64) (VisitDetailOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return VisitDetailOutput{}, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("visit.get.tx_start_error", "err", txErr)
		return VisitDetailOutput{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("visit.get.tx_rollback_error", "err", rbErr)
		}
	}()

	visitWithListing, err := s.visitRepo.GetVisitWithListingByID(ctx, tx, visitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return VisitDetailOutput{}, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.get.get_visit_error", "visit_id", visitID, "err", err)
		return VisitDetailOutput{}, utils.InternalError("")
	}

	owner := s.decorateParticipantSnapshot(ctx, visitWithListing.Owner)
	realtor := s.decorateParticipantSnapshot(ctx, visitWithListing.Realtor)

	return VisitDetailOutput{
		Visit:      visitWithListing.Visit,
		Listing:    visitWithListing.Listing,
		Owner:      owner,
		Realtor:    realtor,
		Timeline:   buildVisitTimeline(visitWithListing.Visit),
		LiveStatus: computeLiveStatus(visitWithListing.Visit, time.Now().UTC()),
	}, nil
}
