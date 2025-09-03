package configadapter

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	policyport "github.com/giulio-alfieri/toq_server/internal/core/port/policy"
	"gopkg.in/yaml.v3"
)

// FileRuleSource loads transition rules from a YAML file path.
type FileRuleSource struct {
	Path string
}

// Load implements TransitionRuleSource by reading YAML from disk and converting to typed rules.
func (frs *FileRuleSource) Load(ctx context.Context) ([]policyport.TransitionRule, error) {
	_ = ctx // não utilizado, mantido para simetria da interface
	if frs.Path == "" {
		return nil, fmt.Errorf("empty rules file path")
	}
	b, err := os.ReadFile(frs.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}
	var raw []struct {
		Role         string `yaml:"role"`
		From         int    `yaml:"from"`
		Action       int    `yaml:"action"`
		To           int    `yaml:"to"`
		Notification int    `yaml:"notification"`
		Priority     int    `yaml:"priority"`
	}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules yaml: %w", err)
	}
	rules := make([]policyport.TransitionRule, 0, len(raw))
	for i, r := range raw {
		tr := policyport.TransitionRule{
			Role:         permissionmodel.RoleSlug(r.Role),
			From:         permissionmodel.UserRoleStatus(r.From),
			Action:       usermodel.ActionFinished(r.Action),
			To:           permissionmodel.UserRoleStatus(r.To),
			Notification: globalmodel.NotificationType(r.Notification),
			Priority:     r.Priority,
		}
		// validação leve
		if !tr.Role.IsValid() {
			slog.Error("Invalid role in status rule", "index", i, "role", r.Role)
			return nil, fmt.Errorf("invalid role at index %d", i)
		}
		rules = append(rules, tr)
	}
	return rules, nil
}
