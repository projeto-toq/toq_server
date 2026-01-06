package proposalservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CancelProposal lets the realtor withdraw a pending proposal.
func (s *proposalService) CancelProposal(ctx context.Context, input StatusChangeInput) (err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ProposalID <= 0 {
		return derrors.Validation("proposalId must be greater than zero", map[string]any{"proposalId": "required"})
	}
	if input.Actor.UserID <= 0 {
		return derrors.Auth("actor metadata missing")
	}
	if input.Actor.RoleSlug != permissionmodel.RoleSlugRealtor {
		return derrors.Forbidden("only realtors can cancel proposals")
	}

	var tx *sql.Tx
	tx, err = s.globalSvc.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.cancel.tx_start_error", "err", err, "proposal_id", input.ProposalID)
		return derrors.Infra("failed to start transaction", err)
	}
	defer s.rollbackOnError(ctx, tx, &err)

	var proposal proposalmodel.ProposalInterface
	proposal, err = s.proposalRepo.GetProposalByIDForUpdate(ctx, tx, input.ProposalID)
	if err != nil {
		return s.mapProposalError(err)
	}

	if proposal.RealtorID() != input.Actor.UserID {
		logger.Warn("proposal.cancel.unauthorized_actor", "proposal_id", input.ProposalID, "actor_id", input.Actor.UserID)
		return derrors.Forbidden("only the author can cancel the proposal")
	}
	if proposal.Status() != proposalmodel.StatusPending {
		return derrors.Conflict("only pending proposals can be cancelled", nil)
	}

	now := time.Now().UTC()
	proposal.SetStatus(proposalmodel.StatusCancelled)
	proposal.SetCancelledAt(sql.NullTime{Valid: true, Time: now})
	proposal.SetAcceptedAt(sql.NullTime{})
	proposal.SetRejectedAt(sql.NullTime{})
	proposal.SetRejectionReason(sql.NullString{})

	if err = s.proposalRepo.UpdateProposalStatus(ctx, tx, proposal, proposalmodel.StatusPending); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return derrors.Conflict("proposal status changed while cancelling", nil)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.cancel.persist_error", "err", err, "proposal_id", input.ProposalID)
		return derrors.Infra("failed to cancel proposal", err)
	}

	if err = s.listingRepo.UpdateProposalFlags(ctx, tx, listingrepository.ProposalFlagsUpdate{
		ListingIdentityID:  proposal.ListingIdentityID(),
		HasPending:         false,
		HasAccepted:        false,
		AcceptedProposalID: sql.NullInt64{},
	}); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.cancel.flags_error", "err", err, "listing_identity_id", proposal.ListingIdentityID())
		return derrors.Infra("failed to update listing proposal flags", err)
	}

	auditMsg := fmt.Sprintf("proposal_cancelled:%d", proposal.ID())
	if err = s.globalSvc.CreateAudit(ctx, tx, globalmodel.TableProposals, auditMsg, input.Actor.UserID); err != nil {
		return err
	}

	if err = s.globalSvc.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.cancel.commit_error", "err", err, "proposal_id", proposal.ID())
		return derrors.Infra("failed to commit proposal cancellation", err)
	}

	logger.Info("proposal.cancel.success",
		"proposal_id", proposal.ID(),
		"listing_identity_id", proposal.ListingIdentityID(),
		"realtor_id", proposal.RealtorID(),
	)

	go s.notifyProposalStatusChange(context.Background(), proposal, proposal.OwnerID(), "proposal_cancelled", "A proposta foi cancelada pelo corretor.")

	return nil
}
