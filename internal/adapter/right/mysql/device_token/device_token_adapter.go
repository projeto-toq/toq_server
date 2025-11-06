package mysqldevicetokenadapter

import (
	mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
)

// DeviceTokenAdapter implements DeviceTokenRepoPortInterface using MySQL
type DeviceTokenAdapter struct {
	db      *mysqladapter.Database
	metrics metricsport.MetricsPortInterface
}

// NewDeviceTokenAdapter creates a new DeviceTokenAdapter
func NewDeviceTokenAdapter(
	db *mysqladapter.Database,
	metrics metricsport.MetricsPortInterface,
) *DeviceTokenAdapter {
	return &DeviceTokenAdapter{
		db:      db,
		metrics: metrics,
	}
}
