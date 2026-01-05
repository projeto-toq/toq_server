package mysqlproposaladapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	proposal_repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/proposal_repository"
)

// Compile-time check to ensure ProposalAdapter satisfies the repository port.
var _ proposal_repository.Repository = (*ProposalAdapter)(nil)

// ProposalAdapter provides persistence access for proposals using MySQL with instrumentation.
type ProposalAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewProposalAdapter builds a proposal adapter wired with metrics instrumentation.
func NewProposalAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *ProposalAdapter {
	return &ProposalAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}
