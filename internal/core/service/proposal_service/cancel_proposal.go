package proposalservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
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
	prevStatus := string(proposal.Status())

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

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		input.Actor.UserID,
		auditmodel.AuditTarget{Type: auditmodel.TargetProposal, ID: proposal.ID()},
		auditmodel.OperationProposalCancel,
		map[string]any{
			"proposal_id":         proposal.ID(),
			"listing_identity_id": proposal.ListingIdentityID(),
			"owner_id":            proposal.OwnerID(),
			"realtor_id":          proposal.RealtorID(),
			"actor_role":          string(permissionmodel.RoleSlugRealtor),
			"status_from":         prevStatus,
			"status_to":           string(proposal.Status()),
		},
	)

	if err = s.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.cancel.audit_error", "err", err, "proposal_id", proposal.ID())
		return derrors.Infra("failed to record proposal audit", err)
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
