package datasources

import (
	"context"
	"time"
)

type CallbackWorker struct {
	Function func() error
	Period   time.Duration
}

func NewCallbackWorker(function func() error, period time.Duration) *CallbackWorker {
	return &CallbackWorker{Function: function, Period: period}
}

func (c *CallbackWorker) Run(ctx context.Context, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.Period):
			err := c.Function()
			if err != nil {
				errorChan <- err
				return
			}
		}
	}
}
