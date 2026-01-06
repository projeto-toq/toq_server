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

// AcceptProposal confirms the owner's approval and updates listing flags.
func (s *proposalService) AcceptProposal(ctx context.Context, input StatusChangeInput) (proposal proposalmodel.ProposalInterface, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return nil, derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ProposalID <= 0 {
		return nil, derrors.Validation("proposalId must be greater than zero", map[string]any{"proposalId": "required"})
	}
	if input.Actor.UserID <= 0 {
		return nil, derrors.Auth("actor metadata missing")
	}
	if input.Actor.RoleSlug != permissionmodel.RoleSlugOwner {
		return nil, derrors.Forbidden("only owners can accept proposals")
	}

	var tx *sql.Tx
	tx, err = s.globalSvc.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.accept.tx_start_error", "err", err, "proposal_id", input.ProposalID)
		return nil, derrors.Infra("failed to start transaction", err)
	}
	defer s.rollbackOnError(ctx, tx, &err)

	proposal, err = s.proposalRepo.GetProposalByIDForUpdate(ctx, tx, input.ProposalID)
	if err != nil {
		return nil, s.mapProposalError(err)
	}

	if proposal.OwnerID() != input.Actor.UserID {
		logger.Warn("proposal.accept.unauthorized_actor", "proposal_id", input.ProposalID, "actor_id", input.Actor.UserID)
		return nil, derrors.Forbidden("only the listing owner can accept proposals")
	}
	if proposal.Status() != proposalmodel.StatusPending {
		return nil, derrors.Conflict("only pending proposals can be accepted", nil)
	}

	now := time.Now().UTC()
	proposal.SetStatus(proposalmodel.StatusAccepted)
	proposal.SetAcceptedAt(sql.NullTime{Valid: true, Time: now})
	proposal.SetRejectedAt(sql.NullTime{})
	proposal.SetCancelledAt(sql.NullTime{})
	proposal.SetRejectionReason(sql.NullString{})

	if err = s.proposalRepo.UpdateProposalStatus(ctx, tx, proposal, proposalmodel.StatusPending); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, derrors.Conflict("proposal status changed before acceptance", nil)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.accept.persist_error", "err", err, "proposal_id", input.ProposalID)
		return nil, derrors.Infra("failed to accept proposal", err)
	}

	if err = s.listingRepo.UpdateProposalFlags(ctx, tx, listingrepository.ProposalFlagsUpdate{
		ListingIdentityID:  proposal.ListingIdentityID(),
		HasPending:         false,
		HasAccepted:        true,
		AcceptedProposalID: sql.NullInt64{Valid: true, Int64: proposal.ID()},
	}); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.accept.flags_error", "err", err, "listing_identity_id", proposal.ListingIdentityID())
		return nil, derrors.Infra("failed to update listing proposal flags", err)
	}

	auditMsg := fmt.Sprintf("proposal_accepted:%d", proposal.ID())
	if err = s.globalSvc.CreateAudit(ctx, tx, globalmodel.TableProposals, auditMsg, input.Actor.UserID); err != nil {
		return nil, err
	}

	if err = s.globalSvc.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.accept.commit_error", "err", err, "proposal_id", proposal.ID())
		return nil, derrors.Infra("failed to commit proposal acceptance", err)
	}

	logger.Info("proposal.accept.success",
		"proposal_id", proposal.ID(),
		"listing_identity_id", proposal.ListingIdentityID(),
		"owner_id", proposal.OwnerID(),
		"realtor_id", proposal.RealtorID(),
	)

	go s.notifyProposalStatusChange(context.Background(), proposal, proposal.RealtorID(), "proposal_accepted", "Sua proposta foi aceita pelo proprietário.")

	return proposal, nil
}
