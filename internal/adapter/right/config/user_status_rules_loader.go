package configadapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	policyport "github.com/giulio-alfieri/toq_server/internal/core/port/policy"
)

// yamlRule represents a rule in configuration. We keep types strongly-typed in memory.
type yamlRule struct {
	Role         string `yaml:"role"`
	From         string `yaml:"from"`
	Action       string `yaml:"action"`
	To           string `yaml:"to"`
	Notification int    `yaml:"notification"`
	Priority     int    `yaml:"priority"`
}

// UserStatusPolicyFromConfig implements UserStatusPolicy using a provided source.
type UserStatusPolicyFromConfig struct {
	src   policyport.TransitionRuleSource
	rules []policyport.TransitionRule
}

// NewUserStatusPolicyFromConfig creates a policy using the given source.
func NewUserStatusPolicyFromConfig(src policyport.TransitionRuleSource) *UserStatusPolicyFromConfig {
	return &UserStatusPolicyFromConfig{src: src}
}

// Evaluate selects the first matching rule ordered by Priority desc then declaration order.
func (p *UserStatusPolicyFromConfig) Evaluate(ctx context.Context, role permissionmodel.RoleSlug, from permissionmodel.UserRoleStatus, action usermodel.ActionFinished) (permissionmodel.UserRoleStatus, globalmodel.NotificationType, bool, error) {
	if p.rules == nil {
		if err := p.Reload(ctx); err != nil {
			return 0, 0, false, err
		}
	}
	for _, r := range p.rules {
		if r.Role == role && r.From == from && r.Action == action {
			changed := r.To != from
			return r.To, r.Notification, changed, nil
		}
		// wildcard by using StatusDeleted+1 not needed; we opt for explicit rules only for safety
	}
	slog.Error("No matching user status transition rule", "role", role, "from", from, "action", action)
	return 0, 0, false, fmt.Errorf("no matching transition rule")
}

// Rules returns an in-memory copy of rules.
func (p *UserStatusPolicyFromConfig) Rules(context.Context) ([]policyport.TransitionRule, error) {
	out := make([]policyport.TransitionRule, len(p.rules))
	copy(out, p.rules)
	return out, nil
}

// Reload reloads rules from the source.
func (p *UserStatusPolicyFromConfig) Reload(ctx context.Context) error {
	if p.src == nil {
		return errors.New("nil rule source")
	}
	rules, err := p.src.Load(ctx)
	if err != nil {
		return err
	}
	// Ordenar por prioridade desc para permitir override explícito
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].Priority > rules[j].Priority })
	p.rules = rules
	slog.Info("User status policy rules loaded", "count", len(rules))
	return nil
}
