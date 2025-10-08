package utils

import (
	"context"
	"log/slog"

	oteltrace "go.opentelemetry.io/otel/trace"
)

// contextualLoggerKey armazena o logger enriquecido no contexto.
type contextualLoggerKey struct{}

// LoggerFromContext returns a logger enriched with correlation identifiers extracted from context.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	if logger, ok := ctx.Value(contextualLoggerKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}

	attrs := contextLoggerAttrs(ctx)
	if len(attrs) == 0 {
		return slog.Default()
	}

	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}

	return slog.Default().With(args...)
}

// ContextWithLogger stores the contextual logger into the provided context.
func ContextWithLogger(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	logger := LoggerFromContext(ctx)
	return context.WithValue(ctx, contextualLoggerKey{}, logger)
}

// contextLoggerAttrs monta a lista de atributos padrão para correlação de logs.
func contextLoggerAttrs(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0, 3)

	if requestID := GetRequestIDFromContext(ctx); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}

	if spanCtx := oteltrace.SpanFromContext(ctx).SpanContext(); spanCtx.IsValid() {
		attrs = append(attrs,
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return attrs
}
