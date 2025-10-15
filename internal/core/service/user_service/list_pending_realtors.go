package userservices

import (
	"context"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// ListPendingRealtors returns realtors pending manual validation with pagination support.
func (us *userService) ListPendingRealtors(ctx context.Context, page, limit int) (ListPendingRealtorsOutput, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	status := permissionmodel.StatusPendingManual
	roleSlug := permissionmodel.RoleSlugRealtor.String()
	isSystemRole := false
	deleted := false

	result, err := us.ListUsers(ctx, ListUsersInput{
		Page:         page,
		Limit:        limit,
		RoleSlug:     roleSlug,
		RoleStatus:   &status,
		IsSystemRole: &isSystemRole,
		Deleted:      &deleted,
	})
	if err != nil {
		return ListPendingRealtorsOutput{}, err
	}

	return ListPendingRealtorsOutput{
		Realtors: result.Users,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	}, nil
}
