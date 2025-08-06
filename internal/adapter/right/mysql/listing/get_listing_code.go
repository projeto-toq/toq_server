package mysqllistingadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (la *ListingAdapter) GetListingCode(ctx context.Context, tx *sql.Tx) (code uint32, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `UPDATE listing_sequence SET id=LAST_INSERT_ID(id+1);`

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		slog.Error("Error preparing statement on msqllistingadapter/GetListingCode", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		slog.Error("Error executing statement on msqllistingadapter/GetListingCode", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}

	code64, err := result.LastInsertId()
	if err != nil {
		slog.Error("Error getting last insert id on msqllistingadapter/GetListingCode", "error", err)
		err = status.Error(codes.Internal, "Internal server error")
		return
	}

	code = uint32(code64)

	return
}
