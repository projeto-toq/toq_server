package proposalservice

import (
	"context"

	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	listingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/listing_repository"
	ownermetricsrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/owner_metrics_repository"
	proposalrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/proposal_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

// Service exposes the orchestration required by handlers.
type Service interface {
	CreateProposal(ctx context.Context, input CreateProposalInput) (proposalmodel.ProposalInterface, error)
	UpdateProposal(ctx context.Context, input UpdateProposalInput) (proposalmodel.ProposalInterface, error)
	CancelProposal(ctx context.Context, input StatusChangeInput) error
	AcceptProposal(ctx context.Context, input StatusChangeInput) (proposalmodel.ProposalInterface, error)
	RejectProposal(ctx context.Context, input StatusChangeInput) (proposalmodel.ProposalInterface, error)
	ListRealtorProposals(ctx context.Context, filter ListFilter) (ListResult, error)
	ListOwnerProposals(ctx context.Context, filter ListFilter) (ListResult, error)
	GetProposalDetail(ctx context.Context, input DetailInput) (DetailResult, error)
}

type proposalService struct {
	proposalRepo proposalrepository.Repository
	listingRepo  listingrepository.ListingRepoPortInterface
	ownerMetrics ownermetricsrepository.Repository
	globalSvc    globalservice.GlobalServiceInterface
	notifier     globalservice.UnifiedNotificationService
	maxDocBytes  int64
}

const defaultMaxDocBytes = 1_000_000

// New builds a Service respecting the factory order (Seção 4 do guia).
func New(
	proposalRepo proposalrepository.Repository,
	listingRepo listingrepository.ListingRepoPortInterface,
	ownerMetrics ownermetricsrepository.Repository,
	globalSvc globalservice.GlobalServiceInterface,
) Service {
	var notifier globalservice.UnifiedNotificationService
	if globalSvc != nil {
		notifier = globalSvc.GetUnifiedNotificationService()
	}

	return &proposalService{
		proposalRepo: proposalRepo,
		listingRepo:  listingRepo,
		ownerMetrics: ownerMetrics,
		globalSvc:    globalSvc,
		notifier:     notifier,
		maxDocBytes:  defaultMaxDocBytes,
	}
}

var _ Service = (*proposalService)(nil)
