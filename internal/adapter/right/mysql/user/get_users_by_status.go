package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUsersByStatus(ctx context.Context, tx *sql.Tx, userRoleStatus permissionmodel.UserRoleStatus, roleSlug permissionmodel.RoleSlug) (users []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// TODO: Implementar busca por status após migração completa do sistema
	// A tabela user_roles mudou de estrutura e não tem mais campos status/role diretos
	// Por enquanto, retornar lista vazia
	slog.Warn("GetUsersByStatus temporarily disabled during migration", "status", userRoleStatus, "role", roleSlug)
	return []usermodel.UserInterface{}, nil

	/* Código original comentado durante migração: mantido somente como referência sem dependências */

}
