package middlewares

import (
	"context"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

var (
	meter            = otel.Meter("toq_server")
	requestCounter   = mustCreateInt64Counter(meter, "toq_server_requests_total")
	errorCounter     = mustCreateInt64Counter(meter, "toq_server_errors_total")
	durationRecorder = mustCreateFloat64Histogram(meter, "toq_server_request_duration_seconds")
)

func mustCreateInt64Counter(meter metric.Meter, name string) metric.Int64Counter {
	counter, err := meter.Int64Counter(name)
	if err != nil {
		panic(err)
	}
	return counter
}

func mustCreateFloat64Histogram(meter metric.Meter, name string) metric.Float64Histogram {
	histogram, err := meter.Float64Histogram(name)
	if err != nil {
		panic(err)
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

		// Start a new span
		tracer := otel.Tracer("toq_server")
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		// Record the start time
		startTime := time.Now()

		// create a new request ID
		requestID := uuid.New().String()

		ctx = context.WithValue(ctx, globalmodel.RequestIDKey, requestID)

		slog.Info("Request received:", "Method:", info.FullMethod)

		// Handle the request
		resp, err := handler(ctx, req)

		// Record the duration
		duration := time.Since(startTime).Seconds()
		durationRecorder.Record(ctx, duration, metric.WithAttributes(attribute.String("method", info.FullMethod)))

		// Increment the request counter
		requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("method", info.FullMethod)))

		// Increment the error counter if there was an error
		if err != nil {
			errorCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("method", info.FullMethod)))
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		return resp, err
	}
}
