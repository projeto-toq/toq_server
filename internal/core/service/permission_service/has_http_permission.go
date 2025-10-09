package permissionservice

import (
	"context"
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HasHTTPPermission verifica se o usuário tem permissão para um endpoint HTTP específico
func (p *permissionServiceImpl) HasHTTPPermission(ctx context.Context, userID int64, method, path string) (allowed bool, err error) {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return false, utils.BadRequest("invalid user id")
	}

	if method == "" || path == "" {
		return false, utils.BadRequest("invalid http method or path")
	}

	logger.Debug("permission.http.check.start", "user_id", userID, "method", method, "path", path)

	// Buscar informações do usuário para obter UserRoleID e RoleStatus usando transação
	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		logger.Error("permission.http.tx_start_failed", "user_id", userID, "error", txErr)
		utils.SetSpanError(ctx, txErr)
		return false, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.http.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	userRole, err := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.http.get_active_role_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}

	if userRole == nil {
		logger.Warn("permission.http.no_active_role", "user_id", userID)
		if cerr := p.globalService.CommitTransaction(ctx, tx); cerr != nil {
			logger.Error("permission.http.tx_commit_failed", "user_id", userID, "error", cerr)
			utils.SetSpanError(ctx, cerr)
			return false, utils.InternalError("")
		}
		return false, nil
	}

	// Mapear HTTP method+path para resource+action
	resource := "http"
	action := fmt.Sprintf("%s:%s", method, path)

	// Criar contexto com as informações completas do usuário
	permContext := permissionmodel.NewPermissionContext(userID, userRole.GetID())
	permContext.AddMetadata("http_method", method)
	permContext.AddMetadata("http_path", path)

	// Usar o método HasPermission para verificar (fora da consulta de role)
	if cerr := p.globalService.CommitTransaction(ctx, tx); cerr != nil {
		logger.Error("permission.http.tx_commit_failed", "user_id", userID, "error", cerr)
		utils.SetSpanError(ctx, cerr)
		return false, utils.InternalError("")
	}
	return p.HasPermission(ctx, userID, resource, action, permContext)
}
