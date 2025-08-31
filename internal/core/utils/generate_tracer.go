package utils

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GenerateTracer creates a new span for business operations with clear naming
// This function is used for internal spans (services, repositories, etc.)
// The span name should be provided explicitly for clarity (e.g., "UserService.GetUserByID")
func GenerateTracer(ctx context.Context) (newctx context.Context, end func(), err error) {
	// Get caller information for debugging purposes only
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		slog.Error("Failed to get caller information")
		err = ErrInternalServer
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		slog.Error("Failed to get function information")
		err = ErrInternalServer
		return
	}

	// Extract package and function name from caller
	fullName := fn.Name()
	parts := strings.Split(fullName, "/")
	lastPart := parts[len(parts)-1]
	nameParts := strings.Split(lastPart, ".")

	if len(nameParts) < 2 {
		slog.Error("Failed to get package and function name", "full_name", fullName)
		err = ErrInternalServer
		return
	}

	packageName := strings.Join(nameParts[:len(nameParts)-1], ".")
	functionName := nameParts[len(nameParts)-1]

	// Create operation name for business logic
	// Format: "PackageName.FunctionName" (e.g., "UserService.GetUserByID")
	operationName := fmt.Sprintf("%s.%s", packageName, functionName)

	// Get request ID for correlation (optional)
	requestID := GetRequestIDFromContext(ctx)
	if requestID == "" {
		// Continue without request ID - internal operations might not have HTTP context
		slog.Debug("No request ID in context for internal operation", "operation", operationName)
	}

	// Create tracer and span
	tracer := otel.Tracer("toq_server")
	newctx, span := tracer.Start(ctx, operationName)

	// Set attributes for debugging and correlation
	span.SetAttributes(
		attribute.String("code.function", fullName),
		attribute.String("code.namespace", "toq_server"),
		attribute.String("code.filepath", file),
		attribute.Int("code.lineno", line),
	)

	// Add request ID if available for correlation
	if requestID != "" {
		span.SetAttributes(attribute.String("app.request_id", requestID))
	}

	// Create end function with proper error handling
	end = func() {
		// Safe span cleanup with panic recovery
		defer func() {
			if r := recover(); r != nil {
				slog.Warn("Recovered from panic in span.End()",
					"panic", r,
					"operation", operationName,
					"request_id", requestID)
			}
		}()

		// Only proceed if span is not nil
		if span != nil {
			span.End()
		}
	}

	return newctx, end, nil
}

// GenerateBusinessTracer creates a span with explicit operation name for better clarity
// Use this for well-defined business operations where you want to control the span name
func GenerateBusinessTracer(ctx context.Context, operationName string) (context.Context, func(), error) {
	// Get caller information for debugging
	pc, file, line, ok := runtime.Caller(1)
	var fullFunctionName string
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			fullFunctionName = fn.Name()
		}
	}

	// Get request ID for correlation
	requestID := GetRequestIDFromContext(ctx)

	// Create tracer and span with explicit operation name
	tracer := otel.Tracer("toq_server")
	newCtx, span := tracer.Start(ctx, operationName)

	// Set debugging and correlation attributes
	span.SetAttributes(
		attribute.String("app.service", "toq_server"),
		attribute.String("business.operation", operationName),
	)

	if fullFunctionName != "" {
		span.SetAttributes(
			attribute.String("code.function", fullFunctionName),
			attribute.String("code.filepath", file),
			attribute.Int("code.lineno", line),
		)
	}

	if requestID != "" {
		span.SetAttributes(attribute.String("app.request_id", requestID))
	}

	// Create end function with error handling
	end := func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Warn("Recovered from panic in business span.End()",
					"panic", r,
					"operation", operationName,
					"request_id", requestID)
			}
		}()

		if span != nil {
			span.End()
		}
	}

	return newCtx, end, nil
}

// SetSpanError sets error information on the current span
func SetSpanError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return
	}

	// Set error status
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Bool("error", true))

	// Extract HTTP status code from HTTPError if possible
	if httpErr, ok := err.(*HTTPError); ok {
		span.SetAttributes(
			attribute.Int("error.code", httpErr.Code),
			attribute.String("error.message", httpErr.Message),
		)
	} else {
		span.SetAttributes(attribute.String("error.message", err.Error()))
	}
}
