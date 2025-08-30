package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
)

// MetricsMiddleware coleta métricas HTTP para análise no Grafana
// Integra com sistema de métricas via port interface (arquitetura hexagonal)
func MetricsMiddleware(metricsAdapter metricsport.MetricsPortInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Incrementar requests em progresso
		metricsAdapter.SetHTTPRequestsInFlight(1) // Simplified - in production would track actual count

		// Processar request
		c.Next()

		// Coletar métricas após processamento
		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())
		size := int64(c.Writer.Size())

		// Registrar métricas
		metricsAdapter.IncrementHTTPRequests(method, path, status)
		metricsAdapter.ObserveHTTPDuration(method, path, duration)
		metricsAdapter.ObserveHTTPResponseSize(method, path, size)

		// Decrementar requests em progresso
		metricsAdapter.SetHTTPRequestsInFlight(0) // Simplified - in production would track actual count
	})
}
