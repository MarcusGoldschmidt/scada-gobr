package datasources

import (
	"context"
	"github.com/google/uuid"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/shared"
	"time"
)

type CallbackDataSource struct {
	function func() error
	id       shared.CommonId
	enable   bool
}

func NewCallbackDataSource(callback func() error, enable bool) *CallbackDataSource {
	return &CallbackDataSource{callback, uuid.New(), enable}
}

func (c CallbackDataSource) Id() shared.CommonId {
	return c.id
}

func (c CallbackDataSource) Name() string {
	return c.id.String()
}

func (c CallbackDataSource) IsEnable() bool {
	return c.enable
}

func (c CallbackDataSource) GetDataPoints() []Datapoint {
	return []Datapoint{}
}

func (c CallbackDataSource) CreateRuntime(ctx context.Context, p persistence.DataPointPersistence) (DataSourceRuntime, error) {
	runtime := CallbackDataSourceRuntime{uuid.New(), c, c.function, p}

	return runtime, nil
}

type CallbackDataSourceRuntime struct {
	id         shared.CommonId
	dataSource DataSource
	function   func() error

	persistence.DataPointPersistence
}

func (c CallbackDataSourceRuntime) GetDataSource() DataSource {
	return c.dataSource
}

func (c CallbackDataSourceRuntime) Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error {
	defer func() {
		shutdownCompleteChan <- c.id
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(1 * time.Second):
			err := c.function()
			if err != nil {
				return err
			}
		}
	}
}
