package shared

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"sync"
	"time"
)

type TimeMutex struct {
	startWait time.Time
	startLock time.Time
	endLock   time.Time

	mutex sync.Mutex

	logger logger.Logger
}

func NewTimeMutex(logger logger.Logger) *TimeMutex {
	return &TimeMutex{logger: logger}
}

func (t *TimeMutex) Lock() {
	t.startWait = time.Now()
	t.logger.Tracef("waiting for lock")
	t.mutex.Lock()
	t.startLock = time.Now()
	t.logger.Tracef("got lock took [%s]", t.startLock.Sub(t.startWait))
}

func (t *TimeMutex) Unlock() {
	t.mutex.Unlock()
	t.endLock = time.Now()
	t.logger.Tracef("released lock took [%s]", t.startLock.Sub(t.endLock))
}
