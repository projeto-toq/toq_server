package globalservice

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (gs *globalService) GetPrivilegeForCache(ctx context.Context, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (privilege usermodel.PrivilegeInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		return
	}

	privilege, err = gs.getPrivilegeForCache(ctx, tx, service, method, role)
	if err != nil {
		gs.RollbackTransaction(ctx, tx)
		return
	}

	err = gs.CommitTransaction(ctx, tx)
	if err != nil {
		gs.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (gs *globalService) getPrivilegeForCache(ctx context.Context, tx *sql.Tx, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (privilege usermodel.PrivilegeInterface, err error) {
	privilege, err = gs.globalRepo.LoadGRPCAccess(ctx, tx, service, method, role)
	if err != nil {
		if codes.NotFound == status.Code(err) {
			return nil, nil
		} else {
			return
		}
	}

	return
}
