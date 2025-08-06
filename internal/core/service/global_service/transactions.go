package globalservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (gs *globalService) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err = gs.globalRepo.StartTransaction(ctx)
	if err != nil {
		slog.Error("Error starting transaction", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return
}

func (gs *globalService) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	err = gs.globalRepo.RollbackTransaction(ctx, tx)
	if err != nil {
		slog.Error("Error rolling back transaction", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	return
}

func (gs *globalService) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	err = gs.globalRepo.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("Error committing transaction", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	return
}
