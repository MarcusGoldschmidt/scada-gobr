package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/util"
)

func (s *Scadagobr) Trace(ctx context.Context, action string, levels ...logger.LogLevel) (context.Context, *util.TraceElapsed) {
	return util.NewTraceElapsed(ctx, s.Logger, s.TimeProvider, action, levels...)
}
