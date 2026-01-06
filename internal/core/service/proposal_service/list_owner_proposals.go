package proposalservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ListOwnerProposals lists proposals received by the owner.
func (s *proposalService) ListOwnerProposals(ctx context.Context, filter ListFilter) (ListResult, error) {
	if filter.Actor.RoleSlug != permissionmodel.RoleSlugOwner {
		return ListResult{}, derrors.Forbidden("only owners can view received proposals")
	}
	return s.listProposals(ctx, proposalmodel.ActorScopeOwner, filter)
}
