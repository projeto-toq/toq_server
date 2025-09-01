package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		return nil, err
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	configuration = make(map[string]string)

	for _, entity := range entities {
		key, ok := entity[1].([]byte)
		if !ok {
			slog.Error("mysqlglobaladapter/GetConfiguration: error converting key to []byte", "key", entity[1])
			return nil, errors.New("configuration key conversion failed")
		}

		value, ok := entity[2].([]byte)
		if !ok {
			slog.Error("mysqlglobaladapter/GetConfiguration: error converting value to []byte", "value", entity[2])
			return nil, errors.New("configuration value conversion failed")
		}
		configuration[string(key)] = string(value)
	}

	return
}
