package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
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
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid metrics handler")
	}
}
