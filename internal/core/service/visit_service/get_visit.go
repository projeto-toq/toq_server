package visitservice

import (
	"context"
	"database/sql"

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

	visit, err := s.visitRepo.GetVisitByID(ctx, tx, visitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return VisitDetailOutput{}, utils.NotFoundError("Visit")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("visit.get.get_visit_error", "visit_id", visitID, "err", err)
		return VisitDetailOutput{}, utils.InternalError("")
	}

	listing, err := s.fetchListingVersionForVisit(ctx, tx, visit)
	if err != nil {
		return VisitDetailOutput{}, err
	}

	return VisitDetailOutput{Visit: visit, Listing: listing}, nil
}
