package datasources

import (
	"context"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/shared"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()

	testLogger := logger.NewTestLogger(t)

	common := NewDataSourceRuntimeManagerCommon(shared.NewCommonId(), "teste", testLogger)

	pointCommon := NewDataPointCommon(shared.NewCommonId(), "randon", true)
	dataPoint := NewRandonValueDataPoint(pointCommon, 0, 100)

	memoryPersistence := persistence.NewInMemoryPersistence()

	worker := NewRandomValueWorker(shared.NewCommonId(), 1*time.Second, []*RandonValueDataPoint{dataPoint}, memoryPersistence)

	common.WithWorker(worker)
	err := common.Run(ctx)
	if err != nil {
		t.Error(err)
	}
	<-time.After(5 * time.Second)
	err = common.Stop(ctx)
	if err != nil {
		t.Error(err)
	}

	<-time.After(2 * time.Second)

	values, err := memoryPersistence.GetPointValues(worker.dataSourceId)
	if err != nil {
		t.Error(err)
	}

	if len(values) != 4 {
		t.Errorf("got %d", len(values))
	}
}
