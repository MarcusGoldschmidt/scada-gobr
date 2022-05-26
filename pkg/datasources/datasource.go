package datasources

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/shared"
	"sync"
	"time"
)

type DataSourceStatus int8

const (
	Initial DataSourceStatus = iota
	Running
	Stopping
	Stopped
	Disabled
	Error
)

var ErrWorkerNotFound = errors.New("worker not found")

type Datapoint interface {
	Id() shared.CommonId
	Name() string
}

type DataSourceWorker interface {
	DataSourceId() shared.CommonId
	Work(ctx context.Context, confirmShutdown chan bool, errorChan chan error)
}

//DataSourceWorkerFunc
//
//func(ctx context.Context, confirmShutdown chan uuid.UUID, errorChan chan error) {
//	defer func() {
//		confirmShutdown <- c.runtimeId
//	}()
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		// Period or worker
//		case err := <-work():
//			if err != nil {
//				errorChan <- err
//			}
//		}
//	}
//}()
type DataSourceWorkerFunc func(ctx context.Context, confirmShutdown chan bool, errorChan chan error)

func (f DataSourceWorkerFunc) Work(ctx context.Context, confirmShutdown chan bool, errorChan chan error) {
	f(ctx, confirmShutdown, errorChan)
}

type DataSourceRuntimeManager interface {
	Id() shared.CommonId
	Name() string

	RuntimeId() uuid.UUID
	Status() DataSourceStatus
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
}

type DataSourceRuntimeManagerCommon struct {
	id   shared.CommonId
	name string

	runtimeId       uuid.UUID
	status          DataSourceStatus
	errorReason     error
	worker          DataSourceWorker
	mutex           sync.Mutex
	shutdown        func()
	confirmShutdown chan bool

	logger logger.Logger
}

func (c *DataSourceRuntimeManagerCommon) Id() shared.CommonId {
	return c.id
}

func (c *DataSourceRuntimeManagerCommon) Name() string {
	return c.name
}

func (c *DataSourceRuntimeManagerCommon) RuntimeId() uuid.UUID {
	return c.runtimeId
}

func (c *DataSourceRuntimeManagerCommon) Status() DataSourceStatus {
	return c.status
}

func (c *DataSourceRuntimeManagerCommon) Restart(ctx context.Context) error {
	err := c.Stop(ctx)
	if err != nil {
		return err
	}
	err = c.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewDataSourceRuntimeManagerCommon(id shared.CommonId, name string, logger logger.Logger) *DataSourceRuntimeManagerCommon {
	return &DataSourceRuntimeManagerCommon{
		id:        id,
		name:      name,
		runtimeId: uuid.New(),
		status:    Initial,
		mutex:     sync.Mutex{},
		logger:    logger,
	}
}

func (c *DataSourceRuntimeManagerCommon) WithWorker(worker DataSourceWorker) {
	c.worker = worker
}

func (c *DataSourceRuntimeManagerCommon) Run(ctx context.Context) error {
	if c.worker == nil {
		return ErrWorkerNotFound
	}

	if c.status == Running {
		return fmt.Errorf("datasource %s with runtimeId: %s is aready running", c.id, c.runtimeId.String())
	}

	c.mutex.Lock()

	c.confirmShutdown = make(chan bool)
	ctx, c.shutdown = context.WithCancel(ctx)

	c.status = Running

	c.mutex.Unlock()

	errorChan := make(chan error)

	go c.worker.Work(ctx, c.confirmShutdown, errorChan)

	go func() {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errorChan:
			if ok {
				c.status = Error
				c.errorReason = err
				c.shutdown()
			}
			return
		}
	}()

	return nil
}

func (c *DataSourceRuntimeManagerCommon) Stop(ctx context.Context) error {
	if c.status != Running {
		return fmt.Errorf("datasource %s with runtimeId: %s is not running aready running", c.id, c.runtimeId.String())
	}

	start := time.Now()

	c.logger.Infof("Shutdown datapoint runtime %s", c.runtimeId.String())

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.status = Stopping

	c.shutdown()
	<-c.confirmShutdown
	close(c.confirmShutdown)
	c.confirmShutdown = nil

	c.status = Stopped

	c.logger.Infof("Shutdown completed of data source id: %s take %d Milliseconds", c.runtimeId.String(), time.Since(start).Milliseconds())

	return nil
}

type DataPointCommon struct {
	id       shared.CommonId
	name     string
	isEnable bool
}

func NewDataPointCommon(id shared.CommonId, name string, isEnable bool) *DataPointCommon {
	return &DataPointCommon{id: id, name: name, isEnable: isEnable}
}
