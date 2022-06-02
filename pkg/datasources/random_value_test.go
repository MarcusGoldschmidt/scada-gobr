package datasources

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/in_memory"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()

	testLogger := logger.NewTestLogger(t)

	common := NewDataSourceRuntimeManagerCommon(shared.NewCommonId(), "teste", testLogger)

	pointCommon := NewDataPointCommon(shared.NewCommonId(), "randon", true)
	dataPoint := NewRandomValueDataPoint(pointCommon, 0, 100)

	memoryPersistence := in_memory.NewInMemoryPersistence()

	worker := NewRandomValueWorker(shared.NewCommonId(), 1*time.Second, []*RandomValueDataPoint{dataPoint}, memoryPersistence)

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
