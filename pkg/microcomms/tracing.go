package microcomms

import (
    "context"
    "log"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

// TracingConfig holds configuration for tracing
type TracingConfig struct {
    ServiceName    string
    TracingEnabled bool
}

// InitTracing initializes distributed tracing
func InitTracing(config TracingConfig) {
    // This is a simplified implementation
    // In a real application, you would configure and initialize
    // OpenTelemetry properly with a tracer provider and exporters
    if !config.TracingEnabled {
        return
    }
    
    log.Println("Tracing initialized for service:", config.ServiceName)
}

// StartSpan starts a new tracing span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
    tracer := otel.Tracer("github.com/pramithamj/microcomms")
    return tracer.Start(ctx, name)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attributes ...trace.SpanEndOption) {
    span := trace.SpanFromContext(ctx)
    span.AddEvent(name)
}

// SpanFromContext gets the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
    return trace.SpanFromContext(ctx)
}