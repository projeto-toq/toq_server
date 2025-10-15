package mysqlpermissionadapter

import (
	"context"
	"log/slog"

	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type PermissionAdapter struct {
	db *mysqladapter.Database
}

func NewPermissionAdapter(db *mysqladapter.Database) permissionrepository.PermissionRepositoryInterface {
	return &PermissionAdapter{
		db: db,
	}
}

func startPermissionOperation(ctx context.Context) (context.Context, func(), *slog.Logger, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ctx, func() {}, nil, err
	}
	//TODO está poluiindo o tracer com esta função centralizada
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx).With("component", "mysql.permission")

	return ctx, spanEnd, logger, nil
}
