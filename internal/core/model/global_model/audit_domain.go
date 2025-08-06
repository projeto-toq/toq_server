package globalmodel

import (
	"time"
)

type audit struct {
	id          int64
	executed_at time.Time
	executed_by int64
	table_name  string
	table_id    int64
	action      string
}

func (a *audit) ID() int64 {
	return a.id
}

func (a *audit) SetID(id int64) {
	a.id = id
}

func (a *audit) ExecutedAt() time.Time {
	return a.executed_at
}

func (a *audit) SetExecutedAt(executedAt time.Time) {
	a.executed_at = executedAt
}

func (a *audit) ExecutedBy() int64 {
	return a.executed_by
}

func (a *audit) SetExecutedBy(executedBy int64) {
	a.executed_by = executedBy
}

func (a *audit) TableName() string {
	return a.table_name
}

func (a *audit) SetTableName(tableName TableName) {
	a.table_name = tableName.String()
}

func (a *audit) TableID() int64 {
	return a.table_id
}

func (a *audit) SetTableID(tableID int64) {
	a.table_id = tableID
}

func (a *audit) Action() string {
	return a.action
}

func (a *audit) SetAction(action string) {
	a.action = action
}
