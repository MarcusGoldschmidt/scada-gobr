package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/util"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type TraceElapsed struct {
	start        time.Time
	end          time.Time
	name         string
	logger       logger.Logger
	timeProvider providers.TimeProvider
	ctx          context.Context
	logLevel     logger.LogLevel
	Span         trace.Span
}

func (t *TraceElapsed) PrintF(message string, params ...interface{}) {
	switch t.logLevel {
	case logger.LogTrace:
		t.logger.Tracef(message, params...)
	case logger.LogDebug:
		t.logger.Debugf(message, params...)
	case logger.LogInfo:
		t.logger.Infof(message, params...)
	case logger.LogWarn:
		t.logger.Warningf(message, params...)
	case logger.LogError:
		t.logger.Errorf(message, params...)
	}
}

func (t *TraceElapsed) End() {
	took := t.timeProvider.GetCurrentTime().Sub(t.start).String()
	t.PrintF("%s took [%s]", t.name, took)
	t.Span.End()
}

func (s *Scadagobr) Trace(ctx context.Context, action string, levels ...logger.LogLevel) (context.Context, *TraceElapsed) {
	level := logger.LogInfo

	if len(levels) > 0 {
		level = levels[0]
	}

	ctx, span := util.Tracer.Start(ctx, action)

	traceElapsed := &TraceElapsed{
		timeProvider: s.timeProvider,
		start:        s.timeProvider.GetCurrentTime(),
		logger:       s.Logger,
		name:         action,
		ctx:          ctx,
		logLevel:     level,
		Span:         span,
	}

	traceElapsed.PrintF("Setting up Scadagobr")

	return ctx, traceElapsed
}
