package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ga *GlobalAdapter) GetConfiguration(ctx context.Context, tx *sql.Tx) (configuration map[string]string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM configuration;`

	entities, err := ga.Read(ctx, tx, query)
	if err != nil {
		slog.Error("mysqlglobaladapter/GetConfiguration: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "Error reading configuration from database")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "Configuration not found")
	}

	configuration = make(map[string]string)

	for _, entity := range entities {
		key, ok := entity[1].([]byte)
		if !ok {
			slog.Error("mysqlglobaladapter/GetConfiguration: error converting key to []byte", "key", entity[1])
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		value, ok := entity[2].([]byte)
		if !ok {
			slog.Error("mysqlglobaladapter/GetConfiguration: error converting value to []byte", "value", entity[2])
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		configuration[string(key)] = string(value)
	}

	return
}
