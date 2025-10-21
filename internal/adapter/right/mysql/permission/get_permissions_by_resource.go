package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// GetPermissionsByResource foi descontinuado após a migração para ações HTTP-only.
// Mantido apenas para compatibilidade, retornando erro para indicar a mudança.
func (pa *PermissionAdapter) GetPermissionsByResource(ctx context.Context, tx *sql.Tx, resource string) (permissions []permissionmodel.PermissionInterface, err error) {
	return nil, fmt.Errorf("GetPermissionsByResource is deprecated: permissions agora são identificadas apenas por ação HTTP")
}
