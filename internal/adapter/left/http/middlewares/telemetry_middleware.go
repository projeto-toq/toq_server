package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// TelemetryMiddleware adds OpenTelemetry tracing to HTTP requests
// Follows OpenTelemetry HTTP semantic conventions for span naming and attributes
// Integrates with metrics collection and provides request correlation
func TelemetryMiddleware(metricsAdapter metricsport.MetricsPortInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get request ID for correlation with logs
		requestID := utils.GetRequestIDFromGinContext(c)
		if requestID == "" {
			requestID = "unknown"
		}

		// Create span following OpenTelemetry HTTP semantic conventions
		// Format: "{HTTP_METHOD} {HTTP_ROUTE}" (e.g., "GET /api/users/{id}")
		method := c.Request.Method
		path := c.Request.URL.Path
		spanName := fmt.Sprintf("%s %s", method, path)

		// Create tracer and span
		tracer := otel.Tracer("toq_server")
		ctx, span := tracer.Start(ctx, spanName)

		// Set OpenTelemetry HTTP semantic convention attributes
		span.SetAttributes(
			semconv.HTTPMethodKey.String(method),
			semconv.HTTPRouteKey.String(path),
			semconv.HTTPSchemeKey.String(c.Request.URL.Scheme),
			attribute.String("http.host", c.Request.Host),
			semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
			attribute.String("http.client_ip", c.ClientIP()),
		)

		// Add custom attributes for correlation and debugging
		span.SetAttributes(
			attribute.String("app.request_id", requestID),
			attribute.String("app.service", "toq_server"),
			attribute.String("app.version", "1.0.0"),
		)

		// Add query parameters if present
		if query := c.Request.URL.RawQuery; query != "" {
			span.SetAttributes(attribute.String("http.query_string", query))
		}

		// Update request context with tracing context
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
			if statusCode >= 400 {
				span.SetAttributes(attribute.Bool("error", true))
				if statusCode >= 500 {
					span.SetAttributes(attribute.String("error.type", "server_error"))
				} else {
					span.SetAttributes(attribute.String("error.type", "client_error"))
				}
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
