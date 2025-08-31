package permissionservice

import (
	"log/slog"

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
	if len(conditions) == 0 {
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
		// Verifica se o usuário é dono do recurso através dos metadados
		if resourceOwnerID, exists := context.Metadata["resource_owner_id"]; exists {
			if ownerID, ok := resourceOwnerID.(int64); ok {
				return ownerID == context.UserID
			}
		}
		// Fallback: verifica se o resource_id é igual ao user_id (auto-ownership)
		if resourceID, exists := context.Metadata["resource_id"]; exists {
			if resID, ok := resourceID.(int64); ok {
				return resID == context.UserID
			}
		}
		return false

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
	// Com a nova estrutura, verificamos o UserRoleID e RoleStatus diretamente
	// Esta função pode ser simplificada ou usar lookup de role baseado no UserRoleID

	// Para agora, vamos fazer uma verificação básica se o usuário tem role ativo
	if !context.IsActive() {
		return false
	}

	switch roleData := roles.(type) {
	case string:
		// Role único - verificar se usuário tem role ativo
		return e.hasRoleByID(context.UserRoleID)

	case []interface{}:
		// Lista de roles (OR lógico)
		for _, role := range roleData {
			if _, ok := role.(string); ok {
				if e.hasRoleByID(context.UserRoleID) {
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

// hasRoleByID verifica se o usuário tem um role específico baseado no UserRoleID
func (e *ConditionEvaluator) hasRoleByID(userRoleID int64) bool {
	// Para simplificar por agora, apenas verifica se tem um UserRoleID válido
	// Esta função deveria fazer lookup do role slug baseado no UserRoleID
	// Por enquanto, retorna true se o usuário tem role ativo (UserRoleID > 0)

	// TODO: Implementar lookup real do role baseado no UserRoleID
	// Isso requereria acesso ao repositório ou serviço de role

	return userRoleID > 0
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
