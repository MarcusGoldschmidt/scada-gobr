package util

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"sync"
)

type TraceMutex struct {
	lock         *sync.Mutex
	logger       logger.Logger
	timeProvider providers.TimeProvider
}

func NewTraceMutex(logger logger.Logger, timeProvider providers.TimeProvider) *TraceMutex {
	return &TraceMutex{
		lock:         &sync.Mutex{},
		logger:       logger,
		timeProvider: timeProvider,
	}
}

func (t *TraceMutex) Lock(ctx context.Context, action string) {
	ctx, trace := NewTraceElapsed(ctx, t.logger, t.timeProvider, "Acquiring lock: "+action)
	t.lock.Lock()
	trace.End()
}

func (t *TraceMutex) Unlock() {
	t.lock.Unlock()
}
