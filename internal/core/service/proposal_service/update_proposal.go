package proposalservice

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateProposal allows the author to edit a pending proposal text/document.
func (s *proposalService) UpdateProposal(ctx context.Context, input UpdateProposalInput) (proposal proposalmodel.ProposalInterface, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return nil, derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = s.validateUpdateInput(input); err != nil {
		return nil, err
	}
	if input.Document != nil {
		if err = s.validateDocumentPayload(input.Document); err != nil {
			return nil, err
		}
	}

	var tx *sql.Tx
	tx, err = s.globalSvc.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.update.tx_start_error", "err", err, "proposal_id", input.ProposalID)
		return nil, derrors.Infra("failed to start transaction", err)
	}
	defer s.rollbackOnError(ctx, tx, &err)

	proposal, err = s.proposalRepo.GetProposalByIDForUpdate(ctx, tx, input.ProposalID)
	if err != nil {
		return nil, s.mapProposalError(err)
	}

	if proposal.RealtorID() != input.EditorID {
		logger.Warn("proposal.update.unauthorized_actor", "proposal_id", input.ProposalID, "actor_id", input.EditorID)
		return nil, derrors.Forbidden("only the author can edit the proposal")
	}
	if proposal.Status() != proposalmodel.StatusPending {
		return nil, derrors.Conflict("only pending proposals can be edited", nil)
	}

	trimmed := strings.TrimSpace(input.ProposalText)
	if trimmed == "" {
		return nil, derrors.Validation("proposal text is required", map[string]any{"proposalText": "required"})
	}
	proposal.SetProposalText(trimmed)

	if err = s.proposalRepo.UpdateProposalText(ctx, tx, proposal); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, derrors.Conflict("proposal status changed while editing", nil)
		}
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.update.persist_error", "err", err, "proposal_id", input.ProposalID)
		return nil, derrors.Infra("failed to update proposal text", err)
	}

	if input.Document != nil {
		if err = s.createProposalDocument(ctx, tx, proposal.ID(), input.Document); err != nil {
			return nil, err
		}
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		input.EditorID,
		auditmodel.AuditTarget{Type: auditmodel.TargetProposal, ID: proposal.ID()},
		auditmodel.OperationUpdate,
		map[string]any{
			"proposal_id":         proposal.ID(),
			"listing_identity_id": proposal.ListingIdentityID(),
			"owner_id":            proposal.OwnerID(),
			"realtor_id":          proposal.RealtorID(),
			"actor_role":          string(permissionmodel.RoleSlugRealtor),
			"status_from":         string(proposal.Status()),
			"status_to":           string(proposal.Status()),
			"updated_fields":      []string{"proposal_text"},
			"document_updated":    input.Document != nil,
		},
	)

	if err = s.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.update.audit_error", "err", err, "proposal_id", proposal.ID())
		return nil, derrors.Infra("failed to record proposal audit", err)
	}

	if err = s.globalSvc.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.update.commit_error", "err", err, "proposal_id", proposal.ID())
		return nil, derrors.Infra("failed to commit proposal update", err)
	}

	logger.Info("proposal.update.success",
		"proposal_id", proposal.ID(),
		"listing_identity_id", proposal.ListingIdentityID(),
		"realtor_id", proposal.RealtorID(),
	)

	return proposal, nil
}
