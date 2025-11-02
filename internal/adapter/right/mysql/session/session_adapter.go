package sessionmysqladapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	sessionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/session_repository"
)

// Ensure implementation satisfies port interface
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

type SessionAdapter struct {
	mysqladapter.InstrumentedAdapter
}

func NewSessionAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) sessionrepository.SessionRepoPortInterface {
	return &SessionAdapter{
		InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics),
	}
}
