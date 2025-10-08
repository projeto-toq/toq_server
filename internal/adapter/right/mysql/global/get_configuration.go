package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) GetConfiguration(ctx context.Context, tx *sql.Tx) (configuration map[string]string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT * FROM configuration;`

	entities, err := ga.Read(ctx, tx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.global.get_configuration.read_error", "error", err)
		return nil, err
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	configuration = make(map[string]string)

	for _, entity := range entities {
		key, ok := entity[1].([]byte)
		if !ok {
			err := errors.New("configuration key conversion failed")
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.global.get_configuration.key_conversion_error", "key", entity[1], "error", err)
			return nil, err
		}

		value, ok := entity[2].([]byte)
		if !ok {
			err := errors.New("configuration value conversion failed")
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.global.get_configuration.value_conversion_error", "value", entity[2], "error", err)
			return nil, err
		}
		configuration[string(key)] = string(value)
	}

	return
}
