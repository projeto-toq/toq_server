package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// StructuredLoggingMiddleware provides structured JSON logging with stdout/stderr separation
// Follows Go best practices and Google Style Guide
// Integrates with existing slog system
// Filters out telemetry and monitoring requests to reduce log noise
func StructuredLoggingMiddleware() gin.HandlerFunc {
	// Paths that should not generate access logs (telemetry/health/monitoring)
	skipLoggingPaths := map[string]bool{
		"/metrics": true,
		"/healthz": true,
		"/readyz":  true,
	}

	// Create separate handlers for stdout (INFO/DEBUG) and stderr (WARN/ERROR)
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	})

	stderrHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	})

	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		// Prefer named route when available for stability
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// Skip logging for telemetry/monitoring requests
		if skipLoggingPaths[path] || strings.Contains(userAgent, "Prometheus") {
			c.Next()
			return
		}

		// Get request info before processing
		requestID := utils.GetRequestIDFromGinContext(c)
		clientIP := utils.GetClientIPFromGinContext(c)
		// userAgent already declared above for filtering

		// Process request
		c.Next()

		// Calculate processing time
		duration := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Build base log fields
		fields := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.Int("size", size),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
		}

		// Add trace correlation if available (OpenTelemetry)
		if sc := oteltrace.SpanFromContext(c.Request.Context()).SpanContext(); sc.IsValid() {
			fields = append(fields,
				slog.String("trace_id", sc.TraceID().String()),
				slog.String("span_id", sc.SpanID().String()),
			)
		}

		// Add query parameters if present
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			fields = append(fields, slog.String("query", rawQuery))
		}

		// Add user info if available (authenticated request)
		if userInfo, err := utils.GetUserInfoFromGinContext(c); err == nil {
			fields = append(fields,
				slog.Int64("user_id", userInfo.ID),
				slog.Int64("user_role_id", userInfo.UserRoleID),
				slog.String("role_status", userInfo.RoleStatus.String()),
			)
		}

		// Determine appropriate logger and log level based on HTTP status
		var logger *slog.Logger
		var logLevel slog.Level
		var message string
		includeErrorDetails := false

		switch {
		case status >= 500:
			logger = slog.New(stderrHandler)
			logLevel = slog.LevelError
			message = "HTTP Error"
			includeErrorDetails = true // 5xx: incluir stack/callsite e lista de erros
		case status == http.StatusTooManyRequests || status == http.StatusLocked:
			// 429/423: condições excepcionais tratadas como WARN, sem stack
			logger = slog.New(stderrHandler)
			logLevel = slog.LevelWarn
			message = "HTTP Response"
		case status >= 400:
			// 4xx esperados: INFO em stdout, sem detalhes de erro/stack
			logger = slog.New(stdoutHandler)
			logLevel = slog.LevelInfo
			message = "HTTP Response"
		case status >= 300:
			logger = slog.New(stdoutHandler)
			logLevel = slog.LevelInfo
			message = "HTTP Request"
		default:
			logger = slog.New(stdoutHandler)
			logLevel = slog.LevelInfo
			message = "HTTP Request"
		}

		// Add error information only for 5xx responses to avoid noise in expected 4xx
		if includeErrorDetails && len(c.Errors) > 0 {
			errorMessages := make([]string, len(c.Errors))
			for i, ginErr := range c.Errors {
				errorMessages[i] = ginErr.Error()
				// Se o erro implementar DomainErrorWithSource, registramos a origem do primeiro
				if i == 0 {
					if derr, ok := ginErr.Err.(utils.DomainErrorWithSource); ok {
						fn, file, line := derr.Location()
						if fn != "" {
							fields = append(fields, slog.String("function", fn))
						}
						if file != "" {
							fields = append(fields, slog.String("file", file))
						}
						if line > 0 {
							fields = append(fields, slog.Int("line", line))
						}
						if stack := derr.Stack(); len(stack) > 0 {
							fields = append(fields, slog.Any("stack", stack))
						}
						// Enriquecer com código/mensagem se também implementar DomainError
						if d, ok := ginErr.Err.(interface {
							Code() int
							Message() string
						}); ok {
							fields = append(fields,
								slog.Int("error_code", d.Code()),
								slog.String("error_message", d.Message()),
							)
						}
					}
				}
			}
			fields = append(fields, slog.Any("errors", errorMessages))
		}

		// Log with appropriate handler (stdout/stderr separation)
		logger.LogAttrs(context.Background(), logLevel, message, fields...)
	})
}

// RequestLoggingMiddleware is a lighter version for basic request logging
// Can be used in development or for specific routes
func RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Simple structured log
		slog.Info("Request processed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
			"client_ip", c.ClientIP(),
		)
	})
}

// ErrorLoggingMiddleware logs only errors and warnings
// Useful for production environments with high traffic
// Filters out telemetry and monitoring requests to reduce log noise
func ErrorLoggingMiddleware() gin.HandlerFunc {
	// Paths that should not generate access logs (telemetry/health/monitoring)
	skipLoggingPaths := map[string]bool{
		"/metrics": true,
		"/healthz": true,
		"/readyz":  true,
	}

	stderrHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	})

	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		status := c.Writer.Status()

		// Skip logging for telemetry/monitoring requests (even errors)
		if skipLoggingPaths[path] || strings.Contains(userAgent, "Prometheus") {
			return
		}

		// Only log errors and warnings
		if status >= 400 {
			fields := []slog.Attr{
				slog.String("request_id", utils.GetRequestIDFromGinContext(c)),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Int("status", status),
				slog.Duration("duration", time.Since(start)),
				slog.String("client_ip", utils.GetClientIPFromGinContext(c)),
			}

			// Add user info if available
			if userInfo, err := utils.GetUserInfoFromGinContext(c); err == nil {
				fields = append(fields, slog.Int64("user_id", userInfo.ID))
			}

			// Add errors if present
			if len(c.Errors) > 0 {
				errorMessages := make([]string, len(c.Errors))
				for i, err := range c.Errors {
					errorMessages[i] = err.Error()
				}
				fields = append(fields, slog.Any("errors", errorMessages))
			}

			logLevel := slog.LevelWarn
			if status >= 500 {
				logLevel = slog.LevelError
			}

			logger := slog.New(stderrHandler)
			logger.LogAttrs(context.Background(), logLevel, "HTTP Error", fields...)
		}
	})
}

// LoggingConfig allows customization of logging behavior
type LoggingConfig struct {
	EnableStdoutStderrSeparation bool
	LogLevel                     slog.Level
	AddSource                    bool
	LogOnlyErrors                bool
	IncludeRequestBody           bool
	IncludeResponseBody          bool
}

// ConfigurableLoggingMiddleware creates a logging middleware with custom config
func ConfigurableLoggingMiddleware(config LoggingConfig) gin.HandlerFunc {
	if config.LogOnlyErrors {
		return ErrorLoggingMiddleware()
	}

	if config.EnableStdoutStderrSeparation {
		return StructuredLoggingMiddleware()
	}

	// Default single handler
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     config.LogLevel,
		AddSource: config.AddSource,
	})

	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		fields := []slog.Attr{
			slog.String("request_id", utils.GetRequestIDFromGinContext(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("client_ip", utils.GetClientIPFromGinContext(c)),
		}

		logger := slog.New(handler)
		logger.LogAttrs(context.Background(), slog.LevelInfo, "HTTP Request", fields...)
	})
}
