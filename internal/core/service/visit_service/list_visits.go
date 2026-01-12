package visitservice

import (
	"context"
	"time"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListVisits lists visits for owners or requesters and hydrates the active listing snapshot for each item.
//
// Steps:
//  1. Start a read-only transaction to guarantee consistent reads between visit and listing tables.
//  2. Delegate query construction/pagination to the visit repository, which already joins the active listing version.
//  3. Map repository entries into VisitDetailOutput so handlers can reuse VisitDetailToResponse.
//
// Errors:
//   - Propagates validation/domain errors from repository directly (already sanitized there).
//   - Wraps infrastructure errors with InternalError to preserve HTTP abstraction per Section 7.1 of the guide.
func (s *visitService) ListVisits(ctx context.Context, filter listingmodel.VisitListFilter) (VisitListOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return VisitListOutput{}, err
	}
	defer spanEnd()

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		return VisitListOutput{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
		}
	}()

	result, err := s.visitRepo.ListVisits(ctx, tx, filter)
	if err != nil {
		return VisitListOutput{}, err
	}

	currentTime := time.Now().UTC()
	items := make([]VisitDetailOutput, 0, len(result.Visits))
	for _, entry := range result.Visits {
		decoratedOwner := s.decorateParticipantSnapshot(ctx, entry.Owner)
		decoratedRealtor := s.decorateParticipantSnapshot(ctx, entry.Realtor)
		items = append(items, VisitDetailOutput{
			Visit:      entry.Visit,
			Listing:    entry.Listing,
			Owner:      decoratedOwner,
			Realtor:    decoratedRealtor,
			Timeline:   buildVisitTimeline(entry.Visit),
			LiveStatus: computeLiveStatus(entry.Visit, currentTime),
		})
	}

	return VisitListOutput{
		Items: items,
		Total: result.Total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}
