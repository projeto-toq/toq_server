package globalrepository

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

type GlobalRepoPortInterface interface {
	CreateAudit(ctx context.Context, tx *sql.Tx, audit globalmodel.AuditInterface) (err error)

	GetConfiguration(ctx context.Context, tx *sql.Tx) (configuration map[string]string, err error)

	// Transaction related methods
	StartTransaction(ctx context.Context) (tx *sql.Tx, err error)
	RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error)
	CommitTransaction(ctx context.Context, tx *sql.Tx) (err error)

	// GRPC related methods
	LoadGRPCAccess(ctx context.Context, tx *sql.Tx, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (privilege usermodel.PrivilegeInterface, err error)
}
