package policy

import (
	"context"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// TransitionRule represents a single FSM rule for user status transitions.
// It maps (role, from, action) => (to, notification).
type TransitionRule struct {
	Role         permissionmodel.RoleSlug
	From         permissionmodel.UserRoleStatus
	Action       usermodel.ActionFinished
	To           permissionmodel.UserRoleStatus
	Notification globalmodel.NotificationType
	// Priority can be used for tie-breaking when multiple rules match.
	Priority int
}

// UserStatusPolicy defines the contract for evaluating status transitions based on rules.
// Evaluate must be pure and side-effect free: it only selects the appropriate transition.
type UserStatusPolicy interface {
	// Evaluate decides the next status and optional notification for a given (role, fromStatus, action).
	// Returns: toStatus, notification, changed (to != from), error when no rule matches or input invalid.
	Evaluate(ctx context.Context, role permissionmodel.RoleSlug, from permissionmodel.UserRoleStatus, action usermodel.ActionFinished) (permissionmodel.UserRoleStatus, globalmodel.NotificationType, bool, error)

	// Rules returns the in-memory rule set for observability/debug.
	Rules(ctx context.Context) ([]TransitionRule, error)

	// Reload forces reloading rules from the underlying source if supported.
	Reload(ctx context.Context) error
}

// TransitionRuleSource defines a rule provider, typically a YAML/JSON loader.
type TransitionRuleSource interface {
	// Load returns the full rule set to be used by the policy engine.
	Load(ctx context.Context) ([]TransitionRule, error)
}
