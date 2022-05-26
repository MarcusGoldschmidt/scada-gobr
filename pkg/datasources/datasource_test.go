package datasources

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/shared"
	"testing"
	"time"
)

func TestSimpleDataSourceRuntimeManagerCommon(t *testing.T) {

	log := logger.NewTestLogger(t)

	datasource := NewDataSourceRuntimeManagerCommon(shared.CommonId(uuid.New()), "teste", log)
	ctx := context.Background()
	work := make(chan error)

	go func() {
		work <- nil
		<-time.After(1 * time.Second)
		work <- errors.New("Teste")
	}()

	err := datasource.Run(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestSimpleDataSourceRuntimeError(t *testing.T) {

	log := logger.NewTestLogger(t)

	datasource := NewDataSourceRuntimeManagerCommon(shared.CommonId(uuid.New()), "teste", log)
	ctx := context.Background()
	work := make(chan error)

	go func() {
		work <- nil

		work <- nil
	}()

	err := datasource.Run(ctx)
	if err != nil {
		t.Error(err)
	}

	err = datasource.Stop(ctx)
	if err != nil {
		t.Error(err)
	}

	if datasource.status != Error {
		t.Fail()
	}
}
