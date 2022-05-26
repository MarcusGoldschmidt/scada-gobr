package datasources

import (
	"context"
	"time"
)

type CallbackWorker struct {
	Function func() error
	Enable   bool
}

func (c *CallbackWorker) Run(ctx context.Context, confirmShutdown chan bool, errorChan chan error) {
	defer func() {
		confirmShutdown <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			err := c.Function()
			if err != nil {
				errorChan <- err
				return
			}
		}
	}
}
