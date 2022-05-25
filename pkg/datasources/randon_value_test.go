package datasources

import (
	"context"
	"github.com/google/uuid"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/shared"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()
	duration, _ := time.ParseDuration("2s")
	ctx, _ = context.WithTimeout(ctx, duration)

	common := NewDataSourceCommon(uuid.New(), "a", true)

	duration, _ = time.ParseDuration("1s")

	dataSource := NewRandonValueDataSource(common, duration)

	dataPointCommon := NewDataPointCommon(uuid.New(), "a", true)

	dataPoint := NewRandonValueDataPoint(dataPointCommon, 0, 100)

	dataSource.AddDataPoint(dataPoint)

	memoryP := persistence.NewInMemoryPersistence()

	runtime, err := dataSource.CreateRuntime(ctx, memoryP)
	if err != nil {
		t.Error(err)
		return
	}

	shutdown := make(chan shared.CommonId)

	go func() {
		err := runtime.Run(ctx, shutdown)
		if err != nil {
			t.Error(err)
		}
	}()

	<-shutdown
}
