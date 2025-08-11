package middlewares

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	meter            = otel.Meter("toq_server")
	requestCounter   = mustCreateInt64Counter(meter, "toq_server_requests_total")
	errorCounter     = mustCreateInt64Counter(meter, "toq_server_errors_total")
	durationRecorder = mustCreateFloat64Histogram(meter, "toq_server_request_duration_seconds")
)

func mustCreateInt64Counter(meter metric.Meter, name string) metric.Int64Counter {
	if meter == nil {
		slog.Error("Meter is nil when creating counter", "name", name)
		return nil
	}
	counter, err := meter.Int64Counter(name)
	if err != nil {
		slog.Error("Failed to create Int64Counter", "name", name, "error", err)
		return nil
	}
	return counter
}

func mustCreateFloat64Histogram(meter metric.Meter, name string) metric.Float64Histogram {
	if meter == nil {
		slog.Error("Meter is nil when creating histogram", "name", name)
		return nil
	}
	histogram, err := meter.Float64Histogram(name)
	if err != nil {
		slog.Error("Failed to create Float64Histogram", "name", name, "error", err)
		return nil
	}
	return histogram
}

func TelemetryInterceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Start a new span - check if tracer provider is available
		tracer := otel.Tracer("toq_server")
		if tracer == nil {
			slog.Error("Tracer is nil - OpenTelemetry not properly initialized")
			// Continue without tracing if tracer is nil
			ctx = context.WithValue(ctx, globalmodel.RequestIDKey, uuid.New().String())
			return handler(ctx, req)
		}

		ctx, span := tracer.Start(ctx, info.FullMethod)

		// Safe span cleanup - simples e direto
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		// Record the start time
		startTime := time.Now()

		// create a new request ID
		requestID := uuid.New().String()

		ctx = context.WithValue(ctx, globalmodel.RequestIDKey, requestID)

		slog.Info("Request received:", "Method:", info.FullMethod)

		// Handle the request with panic protection
		var resp interface{}
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					stackTrace := debug.Stack()
					slog.Error("Recovered from panic in handler",
						"panic", r,
						"method", info.FullMethod,
						"request_id", requestID,
						"stack_trace", string(stackTrace))
					// Set error response for panic
					err = status.Errorf(grpcCodes.Internal, "internal server error")
					if span != nil {
						span.SetStatus(codes.Error, "panic in handler")
						span.RecordError(err)
					}
				}
			}()
			resp, err = handler(ctx, req)
		}()

		// Record the duration - protect against nil durationRecorder
		if durationRecorder != nil {
			duration := time.Since(startTime).Seconds()
			durationRecorder.Record(ctx, duration, metric.WithAttributes(attribute.String("method", info.FullMethod)))
		}

		// Increment the request counter - protect against nil requestCounter
		if requestCounter != nil {
			requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("method", info.FullMethod)))
		}

		// Increment the error counter if there was an error - protect against nil errorCounter
		if err != nil && errorCounter != nil {
			errorCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("method", info.FullMethod)))
			if span != nil {
				span.SetStatus(codes.Error, err.Error())
				span.RecordError(err)
			}
		} else if span != nil {
			span.SetStatus(codes.Ok, "OK")
		}

		return resp, err
	}
}
