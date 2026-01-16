package proposalservice

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateProposal orchestrates validation, persistence, flag updates and notifications.
func (s *proposalService) CreateProposal(ctx context.Context, input CreateProposalInput) (proposal proposalmodel.ProposalInterface, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return nil, derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err = s.validateCreateInput(input); err != nil {
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
		logger.Error("proposal.create.tx_start_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, derrors.Infra("failed to start transaction", err)
	}
	defer s.rollbackOnError(ctx, tx, &err)

	var identity listingrepository.ListingIdentityRecord
	identity, err = s.listingRepo.GetListingIdentityByID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.create.identity_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, s.mapListingError(err)
	}

	if identity.UserID == input.RealtorID {
		return nil, derrors.Forbidden("owners cannot send proposals to themselves")
	}
	if identity.Deleted {
		return nil, derrors.Conflict("listing identity is inactive", nil)
	}
	if identity.HasAcceptedProposal {
		return nil, derrors.Conflict("listing already has an accepted proposal", nil)
	}
	if identity.HasPendingProposal {
		return nil, derrors.Conflict("listing already has a pending proposal", nil)
	}

	proposal = proposalmodel.NewProposal()
	proposal.SetListingIdentityID(input.ListingIdentityID)
	proposal.SetRealtorID(input.RealtorID)
	proposal.SetOwnerID(identity.UserID)
	proposal.SetProposalText(strings.TrimSpace(input.ProposalText))
	proposal.SetStatus(proposalmodel.StatusPending)
	proposal.SetCreatedAt(time.Now().UTC())

	if err = s.proposalRepo.CreateProposal(ctx, tx, proposal); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.create.persist_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, derrors.Infra("failed to persist proposal", err)
	}

	if input.Document != nil {
		if err = s.createProposalDocument(ctx, tx, proposal.ID(), input.Document); err != nil {
			return nil, err
		}
	}

	flagInput := listingrepository.ProposalFlagsUpdate{
		ListingIdentityID:  proposal.ListingIdentityID(),
		HasPending:         true,
		HasAccepted:        false,
		AcceptedProposalID: sql.NullInt64{},
	}
	if err = s.listingRepo.UpdateProposalFlags(ctx, tx, flagInput); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.create.flags_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, derrors.Infra("failed to update listing proposal flags", err)
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		input.RealtorID,
		auditmodel.AuditTarget{Type: auditmodel.TargetProposal, ID: proposal.ID()},
		auditmodel.OperationProposalCreate,
		map[string]any{
			"proposal_id":         proposal.ID(),
			"listing_identity_id": proposal.ListingIdentityID(),
			"owner_id":            proposal.OwnerID(),
			"realtor_id":          proposal.RealtorID(),
			"actor_role":          string(permissionmodel.RoleSlugRealtor),
			"status_from":         "",
			"status_to":           string(proposal.Status()),
		},
	)

	if err = s.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.create.audit_error", "err", err, "proposal_id", proposal.ID())
		return nil, derrors.Infra("failed to record proposal audit", err)
	}

	if err = s.globalSvc.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.create.commit_error", "err", err, "proposal_id", proposal.ID())
		return nil, derrors.Infra("failed to commit proposal creation", err)
	}

	logger.Info("proposal.create.success",
		"proposal_id", proposal.ID(),
		"listing_identity_id", proposal.ListingIdentityID(),
		"owner_id", proposal.OwnerID(),
		"realtor_id", proposal.RealtorID(),
	)

	go s.notifyOwnerNewProposal(context.Background(), proposal, identity.Code)

	return proposal, nil
}
