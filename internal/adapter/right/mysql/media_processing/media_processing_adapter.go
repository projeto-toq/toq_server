package mysqlmediaprocessingadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	mediaprocessingrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/mediaprocessingrepository"
)

// MediaProcessingAdapter implementa o repositório de mídia no MySQL.
type MediaProcessingAdapter struct {
	mysqladapter.InstrumentedAdapter
}

// NewMediaProcessingAdapter cria uma nova instância instrumentada.
func NewMediaProcessingAdapter(db *mysqladapter.Database, metrics metricsport.MetricsPortInterface) *MediaProcessingAdapter {
	return &MediaProcessingAdapter{InstrumentedAdapter: mysqladapter.NewInstrumentedAdapter(db, metrics)}
}

var _ mediaprocessingrepository.RepositoryInterface = (*MediaProcessingAdapter)(nil)
