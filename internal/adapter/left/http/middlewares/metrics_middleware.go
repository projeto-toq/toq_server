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
		// Prefer the named route pattern when available for stable label cardinality
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Incrementar requests em progresso de forma thread-safe
		metricsAdapter.IncrementHTTPRequestsInFlight()

		// Garantir que o decremento sempre aconteça, mesmo em caso de panic
		defer func() {
			metricsAdapter.DecrementHTTPRequestsInFlight()
		}()

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
	})
}
