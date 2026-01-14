package proposalservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetProposalDetail returns proposal metadata plus documents (base64 encoded) to the owner or the realtor.
func (s *proposalService) GetProposalDetail(ctx context.Context, input DetailInput) (DetailResult, error) {
	if input.ProposalID <= 0 {
		return DetailResult{}, derrors.Validation("proposalId must be greater than zero", map[string]any{"proposalId": "required"})
	}
	if input.Actor.UserID <= 0 {
		return DetailResult{}, derrors.Auth("actor metadata missing")
	}

	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return DetailResult{}, derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalSvc.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("proposal.detail.tx_start_error", "err", txErr, "proposal_id", input.ProposalID)
		return DetailResult{}, derrors.Infra("failed to start transaction", txErr)
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		if rbErr := s.globalSvc.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("proposal.detail.tx_rollback_error", "err", rbErr)
		}
	}()

	proposal, err := s.proposalRepo.GetProposalByID(ctx, tx, input.ProposalID)
	if err != nil {
		return DetailResult{}, s.mapProposalError(err)
	}

	if !s.actorCanViewProposal(input.Actor, proposal) {
		return DetailResult{}, derrors.Forbidden("actor cannot access this proposal")
	}

	if input.Actor.RoleSlug == permissionmodel.RoleSlugOwner && !proposal.FirstOwnerActionAt().Valid {
		seenAt := time.Now().UTC()
		if err := s.proposalRepo.MarkOwnerFirstView(ctx, tx, proposal.ID(), input.Actor.UserID, seenAt); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("proposal.detail.mark_first_view_error", "err", err, "proposal_id", input.ProposalID, "owner_id", input.Actor.UserID)
			return DetailResult{}, derrors.Infra("failed to mark owner first view", err)
		}
		proposal.SetFirstOwnerActionAt(sql.NullTime{Valid: true, Time: seenAt})
	}

	documents, err := s.proposalRepo.ListDocuments(ctx, tx, proposal.ID(), true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			documents = nil
		} else {
			utils.SetSpanError(ctx, err)
			logger.Error("proposal.detail.documents_error", "err", err, "proposal_id", input.ProposalID)
			return DetailResult{}, derrors.Infra("failed to list proposal documents", err)
		}
	}

	realtorSummary, err := s.loadSingleRealtorSummary(ctx, tx, proposal.RealtorID())
	if err != nil {
		return DetailResult{}, err
	}

	ownerSummary, err := s.loadSingleOwnerSummary(ctx, tx, proposal.OwnerID())
	if err != nil {
		return DetailResult{}, err
	}

	listing, listingErr := s.listingRepo.GetActiveListingVersion(ctx, tx, proposal.ListingIdentityID())
	if listingErr != nil && !errors.Is(listingErr, sql.ErrNoRows) {
		utils.SetSpanError(ctx, listingErr)
		logger.Error("proposal.detail.listing_error", "err", listingErr, "listing_identity_id", proposal.ListingIdentityID())
		return DetailResult{}, derrors.Infra("failed to load listing", listingErr)
	}

	if err := s.globalSvc.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.detail.tx_commit_error", "err", err, "proposal_id", proposal.ID())
		return DetailResult{}, derrors.Infra("failed to commit proposal detail transaction", err)
	}
	committed = true

	return DetailResult{Proposal: proposal, Documents: documents, Realtor: realtorSummary, Owner: ownerSummary, Listing: listing}, nil
}
