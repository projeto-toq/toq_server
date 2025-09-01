package globalconverters

import (
	"errors"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

func AuditEntityToDomain(entity []any) (audit globalmodel.AuditInterface, err error) {
	audit = globalmodel.NewAudit()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "ID", entity[0])
		return nil, errors.New("invalid audit ID type")
	}
	audit.SetID(id)

	executed_at, ok := entity[1].(time.Time)
	if !ok {
		slog.Error("Error converting executed_at to time.Time", "executed_at", entity[1])
		return nil, errors.New("invalid executed_at type")
	}
	audit.SetExecutedAt(executed_at)

	executed_by, ok := entity[2].(int64)
	if !ok {
		slog.Error("Error converting executed_by to int64", "executed_by", entity[2])
		return nil, errors.New("invalid executed_by type")
	}
	audit.SetExecutedBy(executed_by)

	table_name, ok := entity[3].([]byte)
	if !ok {
		slog.Error("Error converting table_name to []byte", "table_name", entity[3])
		return nil, errors.New("invalid table_name type")
	}
	audit.SetTableName(globalmodel.TableName(string(table_name)))

	table_id, ok := entity[4].(int64)
	if !ok {
		slog.Error("Error converting table_id to int64", "table_id", entity[4])
		return nil, errors.New("invalid table_id type")
	}
	audit.SetTableID(table_id)

	action, ok := entity[5].([]byte)
	if !ok {
		slog.Error("Error converting action to []byte", "action", entity[5])
		return nil, errors.New("invalid action type")
	}
	audit.SetAction(string(action))

	return
}
