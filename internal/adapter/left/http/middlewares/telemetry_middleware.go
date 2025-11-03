package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// TelemetryMiddleware adds OpenTelemetry tracing to HTTP requests.
// It follows OpenTelemetry HTTP semantic conventions for span naming and attributes,
// integrates with metrics collection and provides request correlation.
// Telemetry/monitoring endpoints are filtered out to reduce trace noise.
func TelemetryMiddleware(metricsAdapter metricsport.MetricsPortInterface) gin.HandlerFunc {
	// Paths that should not generate traces (telemetry/health/monitoring)
	skipTracingPaths := map[string]bool{
		"/metrics": true,
		"/healthz": true,
		"/readyz":  true,
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		// Prefer the named route pattern when available for cardinality stability
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		path := route
		userAgent := c.Request.UserAgent()

		// Skip tracing AND metrics for telemetry/monitoring requests
		if skipTracingPaths[path] || strings.Contains(userAgent, "Prometheus") {
			// Just continue without telemetry overhead
			c.Next()
			return
		}

		environment := globalmodel.GetRuntimeEnvironment()
		ctx := c.Request.Context()

		// Get request ID for correlation with logs
		requestID := utils.GetRequestIDFromGinContext(c)
		if requestID == "" {
			requestID = "unknown"
		}

		// Create span following OpenTelemetry HTTP semantic conventions
		// Format: "{HTTP_METHOD} {HTTP_ROUTE}" (e.g., "GET /api/users/{id}")
		method := c.Request.Method
		// path already declared above for filtering
		spanName := fmt.Sprintf("%s %s", method, path)
		scheme := c.Request.URL.Scheme
		if scheme == "" {
			if c.Request.TLS != nil {
				scheme = "https"
			} else {
				scheme = "http"
			}
		}

		// Create tracer and span
		tracer := otel.Tracer("toq_server")
		ctx, span := tracer.Start(ctx, spanName)
		ctx = utils.ContextWithLogger(ctx)

		// Set OpenTelemetry HTTP semantic convention attributes
		span.SetAttributes(
			semconv.HTTPMethodKey.String(method),
			semconv.HTTPRouteKey.String(path),
			semconv.HTTPSchemeKey.String(scheme),
			attribute.String("http.host", c.Request.Host),
			semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// Add custom attributes for correlation and debugging
		span.SetAttributes(
			attribute.String("app.request_id", requestID),
			attribute.String("app.service", "toq_server"),
			attribute.String("app.version", globalmodel.AppVersion),
			attribute.String("deployment.environment", environment),
		)
		span.SetAttributes(attribute.String("http.route_template", path))

		// Add query parameters if present
		if c.Request.URL.RawQuery != "" {
			span.SetAttributes(attribute.Bool("http.query_params_present", true))
		}

		// Atualiza o contexto da requisição com o contexto de tracing criado acima
		c.Request = c.Request.WithContext(ctx)

		// Ensure span is ended properly
		defer func() {
			// Set response attributes after request processing
			statusCode := c.Writer.Status()
			responseSize := c.Writer.Size()

			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(statusCode),
				attribute.Int("http.response_content_length", responseSize),
			)

			// Set span status based on HTTP status code
			switch {
			case statusCode >= http.StatusInternalServerError:
				span.SetAttributes(attribute.String("error.type", "server_error"))
				span.SetStatus(codes.Error, http.StatusText(statusCode))
				utils.SetSpanError(ctx, fmt.Errorf("http_status_%d", statusCode))
			case statusCode >= http.StatusBadRequest:
				span.SetAttributes(attribute.String("error.type", "client_error"))
				span.SetStatus(codes.Error, http.StatusText(statusCode))
			default:
				span.SetStatus(codes.Ok, http.StatusText(statusCode))
			}

			span.End()
		}()

		// Apply metrics middleware if adapter is provided
		if metricsAdapter != nil {
			MetricsMiddleware(metricsAdapter)(c)
		} else {
			c.Next()
		}
	})
}
