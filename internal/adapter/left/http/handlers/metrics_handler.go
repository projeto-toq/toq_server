package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
)

// MetricsHandler expõe endpoint /metrics para coleta pelo Prometheus
type MetricsHandler struct {
	metricsAdapter metricsport.MetricsPortInterface
}

// NewMetricsHandler cria uma nova instância do handler de métricas
func NewMetricsHandler(metricsAdapter metricsport.MetricsPortInterface) *MetricsHandler {
	return &MetricsHandler{
		metricsAdapter: metricsAdapter,
	}
}

// GetMetrics expõe as métricas no formato Prometheus
// GET /metrics
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	handler := h.metricsAdapter.GetMetricsHandler()

	// Converter para http.Handler
	if httpHandler, ok := handler.(http.Handler); ok {
		httpHandler.ServeHTTP(c.Writer, c.Request)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid metrics handler",
		})
	}
}
