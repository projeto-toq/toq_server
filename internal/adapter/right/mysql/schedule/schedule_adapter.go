package mysqlscheduleadapter

import mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

// ScheduleAdapter provides DB access to listing agendas.
type ScheduleAdapter struct {
	db *mysqladapter.Database
}

// NewScheduleAdapter builds a new adapter.
func NewScheduleAdapter(db *mysqladapter.Database) *ScheduleAdapter {
	return &ScheduleAdapter{db: db}
}
