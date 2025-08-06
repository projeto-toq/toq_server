package globalmodel

import (
	"time"
)

type AuditInterface interface {
	ID() int64
	SetID(id int64)
	ExecutedAt() time.Time
	SetExecutedAt(executedAt time.Time)
	ExecutedBy() int64
	SetExecutedBy(executedBy int64)
	TableName() string
	SetTableName(tableName TableName)
	TableID() int64
	SetTableID(tableID int64)
	Action() string
	SetAction(action string)
}

func NewAudit() AuditInterface {
	return &audit{}
}
