package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ua *UserAdapter) GetUserPhotoByID(ctx context.Context, tx *sql.Tx, id int64) (photo []byte, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT photo FROM users WHERE id = ?;`

	entities, err := ua.Read(ctx, tx, query, id)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserPhotoByID: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetUserPhotoByID: multiple users found with the same id", "ID", id)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if entities[0][0] != nil {
		var ok bool
		photo, ok = entities[0][0].([]byte)
		if !ok {
			slog.Error("mysqluseradapter/GetUserPhotoByID: error converting photo to string", "photo", entities[0][0])
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	return

}
