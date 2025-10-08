package globalconverters

import (
	"fmt"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

func AuditEntityToDomain(entity []any) (audit globalmodel.AuditInterface, err error) {
	audit = globalmodel.NewAudit()

	id, ok := entity[0].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid audit ID type: %v", entity[0])
	}
	audit.SetID(id)

	executed_at, ok := entity[1].(time.Time)
	if !ok {
		return nil, fmt.Errorf("invalid executed_at type: %v", entity[1])
	}
	audit.SetExecutedAt(executed_at)

	executed_by, ok := entity[2].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid executed_by type: %v", entity[2])
	}
	audit.SetExecutedBy(executed_by)

	table_name, ok := entity[3].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid table_name type: %v", entity[3])
	}
	audit.SetTableName(globalmodel.TableName(string(table_name)))

	table_id, ok := entity[4].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid table_id type: %v", entity[4])
	}
	audit.SetTableID(table_id)

	action, ok := entity[5].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid action type: %v", entity[5])
	}
	audit.SetAction(string(action))

	return
}
