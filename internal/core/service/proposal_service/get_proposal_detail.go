package proposalservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
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

	tx, txErr := s.globalSvc.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("proposal.detail.tx_start_error", "err", txErr, "proposal_id", input.ProposalID)
		return DetailResult{}, derrors.Infra("failed to start read-only transaction", txErr)
	}
	defer func() {
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

	return DetailResult{Proposal: proposal, Documents: documents, Realtor: realtorSummary}, nil
}
