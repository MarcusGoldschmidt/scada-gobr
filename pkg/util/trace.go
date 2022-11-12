package util

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

var Tracer = otel.Tracer("Scadagobr")

func NewContextWithTrace(ctx context.Context) context.Context {
	span := trace.SpanFromContext(ctx)

	newCtx := trace.ContextWithSpan(context.Background(), span)

	return newCtx
}

func TraceParent(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)

	if !spanCtx.IsValid() {
		return ""
	}

	return fmt.Sprintf("00-%s-%s-%s", spanCtx.TraceID().String(), spanCtx.SpanID().String(), spanCtx.TraceFlags().String())
}

func FromTraceParent(ctx context.Context, traceParent string) (context.Context, error) {
	splitResult := strings.Split(traceParent, "-")

	if len(splitResult) != 4 {
		return nil, errors.New("invalid traceparent")
	}

	traceId, err := trace.TraceIDFromHex(splitResult[1])

	if err != nil {
		return nil, err
	}

	spanId, err := trace.SpanIDFromHex(splitResult[2])

	if err != nil {
		return nil, err
	}

	flag, err := hex.DecodeString(splitResult[3])
	if err != nil {
		return nil, err
	}

	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceId,
		SpanID:     spanId,
		TraceFlags: trace.TraceFlags(flag[0]),
	}))

	return ctx, nil
}
