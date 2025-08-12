package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sa *SessionAdapter) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	tx, err = sa.db.DB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("Error starting transaction", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return tx, nil
}

func (sa *SessionAdapter) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Rollback()
	if err != nil {
		slog.Error("Error rolling back transaction", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}
	return nil
}

func (sa *SessionAdapter) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Commit()
	if err != nil {
		slog.Error("Error committing transaction", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}
	return nil
}
