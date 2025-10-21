package permissionservice

import (
	"context"
	"strings"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ConditionEvaluator avalia condições de permissão
type ConditionEvaluator struct{}

// NewConditionEvaluator cria uma nova instância do avaliador
func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{}
}

// Evaluate avalia se as condições são atendidas pelo contexto
func (e *ConditionEvaluator) Evaluate(ctx context.Context, conditions map[string]interface{}, permissionCtx *permissionmodel.PermissionContext) bool {
	if len(conditions) == 0 {
		return true // Sem condições = permitido
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	logger.Debug("permission.condition.evaluate", "conditions", conditions, "context", permissionCtx)

	// Verificar condição de proprietário
	if owner, exists := conditions["owner"]; exists {
		if !e.checkOwnerCondition(ctx, owner, permissionCtx) {
			logger.Debug("permission.condition.owner_failed", "owner", owner)
			return false
		}
	}

	// Verificar condição de role
	if roles, exists := conditions["role"]; exists {
		if !e.checkRoleCondition(ctx, roles, permissionCtx) {
			logger.Debug("permission.condition.role_failed", "roles", roles)
			return false
		}
	}

	// Verificar condição de relacionamento
	if related, exists := conditions["related"]; exists {
		if !e.checkRelatedCondition(ctx, related, permissionCtx) {
			logger.Debug("permission.condition.related_failed", "related", related)
			return false
		}
	}

	logger.Debug("permission.condition.passed")
	return true
}

// checkOwnerCondition verifica condições de proprietário
func (e *ConditionEvaluator) checkOwnerCondition(ctx context.Context, owner interface{}, permissionCtx *permissionmodel.PermissionContext) bool {
	logger := utils.LoggerFromContext(ctx)
	ownerStr, ok := owner.(string)
	if !ok {
		return false
	}

	switch ownerStr {
	case "self":
		// Verifica se o usuário é dono do recurso através dos metadados
		if resourceOwnerID, exists := permissionCtx.Metadata["resource_owner_id"]; exists {
			if ownerID, ok := resourceOwnerID.(int64); ok {
				return ownerID == permissionCtx.UserID
			}
		}
		// Fallback: verifica se o resource_id é igual ao user_id (auto-ownership)
		if resourceID, exists := permissionCtx.Metadata["resource_id"]; exists {
			if resID, ok := resourceID.(int64); ok {
				return resID == permissionCtx.UserID
			}
		}
		return false

	case "listing_owner":
		// Verifica se o usuário é proprietário do listing
		if resourceOwner, exists := permissionCtx.Metadata["listing_owner_id"]; exists {
			if ownerID, ok := resourceOwner.(int64); ok {
				return ownerID == permissionCtx.UserID
			}
		}
		return false

	default:
		logger.Warn("permission.condition.owner_unknown", "condition", ownerStr)
		return false
	}
}

// checkRoleCondition verifica condições de role
func (e *ConditionEvaluator) checkRoleCondition(ctx context.Context, roles interface{}, permissionCtx *permissionmodel.PermissionContext) bool {
	logger := utils.LoggerFromContext(ctx)
	if permissionCtx == nil {
		return false
	}

	roleSlug := permissionCtx.GetRoleSlug()
	if roleSlug == "" {
		if rawSlug, ok := permissionCtx.Metadata["role_slug"]; ok {
			if slugStr, okCast := rawSlug.(string); okCast {
				roleSlug = permissionmodel.RoleSlug(strings.ToLower(slugStr))
			}
		}
	}

	switch roleData := roles.(type) {
	case string:
		return compareRoleSlug(roleSlug, roleData)
	case []interface{}:
		for _, role := range roleData {
			if roleStr, ok := role.(string); ok {
				if compareRoleSlug(roleSlug, roleStr) {
					return true
				}
			}
		}
		return false
	case []string:
		for _, roleStr := range roleData {
			if compareRoleSlug(roleSlug, roleStr) {
				return true
			}
		}
		return false
	default:
		logger.Warn("permission.condition.role_invalid_format", "roles", roles)
		return false
	}
}

func compareRoleSlug(current permissionmodel.RoleSlug, expected string) bool {
	if current == "" || expected == "" {
		return false
	}
	return strings.EqualFold(current.String(), expected)
}

// checkRelatedCondition verifica condições de relacionamento
func (e *ConditionEvaluator) checkRelatedCondition(ctx context.Context, related interface{}, permissionCtx *permissionmodel.PermissionContext) bool {
	logger := utils.LoggerFromContext(ctx)
	relatedStr, ok := related.(string)
	if !ok {
		return false
	}

	switch relatedStr {
	case "owner_or_realtor":
		// Verifica se é dono ou corretor relacionado
		return e.isOwnerOrRelatedRealtor(permissionCtx)

	case "same_agency":
		// Verifica se pertence à mesma agência
		return e.isSameAgency(permissionCtx)

	default:
		logger.Warn("permission.condition.related_unknown", "condition", relatedStr)
		return false
	}
}

// isOwnerOrRelatedRealtor verifica se é proprietário ou corretor relacionado
func (e *ConditionEvaluator) isOwnerOrRelatedRealtor(context *permissionmodel.PermissionContext) bool {
	// Verifica se é o proprietário do recurso através dos metadados
	if resourceOwnerID, exists := context.Metadata["resource_owner_id"]; exists {
		if ownerID, ok := resourceOwnerID.(int64); ok && ownerID == context.UserID {
			return true
		}
	}

	// Verifica se é um corretor relacionado (via agência ou convite)
	if realtorIDs, exists := context.Metadata["related_realtor_ids"]; exists {
		if ids, ok := realtorIDs.([]int64); ok {
			for _, id := range ids {
				if id == context.UserID {
					return true
				}
			}
		}
	}

	return false
}

// isSameAgency verifica se pertence à mesma agência
func (e *ConditionEvaluator) isSameAgency(context *permissionmodel.PermissionContext) bool {
	userAgencyID, userExists := context.Metadata["user_agency_id"]
	resourceAgencyID, resourceExists := context.Metadata["resource_agency_id"]

	if !userExists || !resourceExists {
		return false
	}

	userID, userOk := userAgencyID.(int64)
	resourceID, resourceOk := resourceAgencyID.(int64)

	return userOk && resourceOk && userID == resourceID && userID > 0
}
