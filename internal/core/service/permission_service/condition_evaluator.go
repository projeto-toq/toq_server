package permissionservice

import (
	"log/slog"
	"strings"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// ConditionEvaluator avalia condições de permissão
type ConditionEvaluator struct{}

// NewConditionEvaluator cria uma nova instância do avaliador
func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{}
}

// Evaluate avalia se as condições são atendidas pelo contexto
func (e *ConditionEvaluator) Evaluate(conditions map[string]interface{}, context *permissionmodel.PermissionContext) bool {
	if conditions == nil || len(conditions) == 0 {
		return true // Sem condições = permitido
	}

	slog.Debug("Evaluating conditions", "conditions", conditions, "context", context)

	// Verificar condição de proprietário
	if owner, exists := conditions["owner"]; exists {
		if !e.checkOwnerCondition(owner, context) {
			slog.Debug("Owner condition failed", "owner", owner)
			return false
		}
	}

	// Verificar condição de role
	if roles, exists := conditions["role"]; exists {
		if !e.checkRoleCondition(roles, context) {
			slog.Debug("Role condition failed", "roles", roles)
			return false
		}
	}

	// Verificar condição de relacionamento
	if related, exists := conditions["related"]; exists {
		if !e.checkRelatedCondition(related, context) {
			slog.Debug("Related condition failed", "related", related)
			return false
		}
	}

	slog.Debug("All conditions passed")
	return true
}

// checkOwnerCondition verifica condições de proprietário
func (e *ConditionEvaluator) checkOwnerCondition(owner interface{}, context *permissionmodel.PermissionContext) bool {
	ownerStr, ok := owner.(string)
	if !ok {
		return false
	}

	switch ownerStr {
	case "self":
		// Verifica se o usuário é dono do recurso
		return context.OwnerID != nil && *context.OwnerID == context.UserID

	case "listing_owner":
		// Verifica se o usuário é proprietário do listing
		if resourceOwner, exists := context.Metadata["listing_owner_id"]; exists {
			if ownerID, ok := resourceOwner.(int64); ok {
				return ownerID == context.UserID
			}
		}
		return false

	default:
		slog.Warn("Unknown owner condition", "condition", ownerStr)
		return false
	}
}

// checkRoleCondition verifica condições de role
func (e *ConditionEvaluator) checkRoleCondition(roles interface{}, context *permissionmodel.PermissionContext) bool {
	switch roleData := roles.(type) {
	case string:
		// Role único
		return e.hasRole(roleData, context.UserRoles)

	case []interface{}:
		// Lista de roles (OR lógico)
		for _, role := range roleData {
			if roleStr, ok := role.(string); ok {
				if e.hasRole(roleStr, context.UserRoles) {
					return true
				}
			}
		}
		return false

	default:
		slog.Warn("Invalid role condition format", "roles", roles)
		return false
	}
}

// checkRelatedCondition verifica condições de relacionamento
func (e *ConditionEvaluator) checkRelatedCondition(related interface{}, context *permissionmodel.PermissionContext) bool {
	relatedStr, ok := related.(string)
	if !ok {
		return false
	}

	switch relatedStr {
	case "owner_or_realtor":
		// Verifica se é dono ou corretor relacionado
		return e.isOwnerOrRelatedRealtor(context)

	case "same_agency":
		// Verifica se pertence à mesma agência
		return e.isSameAgency(context)

	default:
		slog.Warn("Unknown related condition", "condition", relatedStr)
		return false
	}
}

// hasRole verifica se o usuário tem um role específico
func (e *ConditionEvaluator) hasRole(targetRole string, userRoles []string) bool {
	for _, role := range userRoles {
		if strings.EqualFold(role, targetRole) {
			return true
		}
	}
	return false
}

// isOwnerOrRelatedRealtor verifica se é proprietário ou corretor relacionado
func (e *ConditionEvaluator) isOwnerOrRelatedRealtor(context *permissionmodel.PermissionContext) bool {
	// Verifica se é o proprietário do recurso
	if context.OwnerID != nil && *context.OwnerID == context.UserID {
		return true
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
