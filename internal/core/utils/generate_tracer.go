package utils

import (
	"context"
	"log/slog"
	"runtime"
	"strings"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GenerateTracer(ctx context.Context) (newctx context.Context, end func(), err error) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		slog.Error("Failed to get caller information")
		err = status.Errorf(codes.Internal, "failed to get caller information")
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		slog.Error("Failed to get function information")
		err = status.Errorf(codes.Internal, "failed to get function information")
		return
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, "/")
	lastPart := parts[len(parts)-1]
	nameParts := strings.Split(lastPart, ".")
	if len(nameParts) < 2 {
		slog.Error("Failed to get package and function name", "full_name", fullName)
		err = status.Errorf(codes.Internal, "failed to get package and function name")
		return
	}

	packageName := strings.Join(nameParts[:len(nameParts)-1], ".")
	functionName := nameParts[len(nameParts)-1]

	requestID, ok := ctx.Value(globalmodel.RequestIDKey).(string)
	if !ok {
		slog.Error("Request ID not found in context", "package", packageName, "function", functionName)
		// Return a no-op function instead of creating span with nil tracer
		end = func() {}
		err = status.Error(codes.Unauthenticated, "")
		return ctx, end, err
	}

	tracer := otel.Tracer(requestID)
	newctx, span := tracer.Start(ctx, packageName+"."+functionName)
	end = func() {
		// Safe span cleanup with panic recovery
		defer func() {
			if r := recover(); r != nil {
				slog.Warn("Recovered from panic in span.End()", "panic", r, "function", packageName+"."+functionName)
			}
		}()

		// Only proceed if span is not nil
		if span != nil {
			if err != nil {
				span.SetAttributes(
					attribute.String("error.message", err.Error()),
					attribute.String("error.code", status.Code(err).String()),
				)
			}
			span.End()
		}
	}

	return

}
