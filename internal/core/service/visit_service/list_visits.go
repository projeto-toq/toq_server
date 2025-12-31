package visitservice

import (
	"context"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *visitService) ListVisits(ctx context.Context, filter listingmodel.VisitListFilter) (listingmodel.VisitListResult, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return listingmodel.VisitListResult{}, err
	}
	defer spanEnd()

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return listingmodel.VisitListResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
		}
	}()

	result, err := s.visitRepo.ListVisits(ctx, tx, filter)
	if err != nil {
		return listingmodel.VisitListResult{}, err
	}

	return result, nil
}
