package datasources

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestSimpleDataSourceRuntimeManagerCommon(t *testing.T) {

	ctx := context.Background()
	log := logger.NewTestLogger(t)

	datasource := NewDataSourceRuntimeManagerCommon(shared.CommonId(uuid.New()), "teste", log)

	datasource.WithWorker(DataSourceWorkerFunc(func(ctx context.Context, errorChan chan<- error) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):

			}
		}
	}))

	err := datasource.Run(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestSimpleDataSourceRuntimeError(t *testing.T) {
	ctx := context.Background()
	log := logger.NewTestLogger(t)

	datasource := NewDataSourceRuntimeManagerCommon(shared.CommonId(uuid.New()), "teste", log)
	datasource.WithWorker(DataSourceWorkerFunc(func(ctx context.Context, errorChan chan<- error) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1):
				errorChan <- errors.New("error test")
			}
		}
	}))

	err := datasource.Run(ctx)
	if err != nil {
		t.Error(err)
	}

	<-time.After(time.Second)

	if datasource.status != Error {
		t.Fail()
	}
}
