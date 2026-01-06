package proposalservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ListRealtorProposals returns paginated history filtered by realtor context.
func (s *proposalService) ListRealtorProposals(ctx context.Context, filter ListFilter) (ListResult, error) {
	if filter.Actor.RoleSlug != permissionmodel.RoleSlugRealtor {
		return ListResult{}, derrors.Forbidden("only realtors can view their proposals")
	}
	return s.listProposals(ctx, proposalmodel.ActorScopeRealtor, filter)
}
