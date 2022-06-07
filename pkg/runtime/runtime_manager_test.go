package runtime

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/in_memory"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()

	rt := NewRuntimeManager(logger.NewSimpleLogger("teste", os.Stdout), in_memory.NewInMemoryPersistence())

	testLogger := logger.NewTestLogger(t)

	sum := 0

	worker := datasources.NewCallbackWorker(func() error {
		sum = sum + 1
		t.Logf("Pass %d seconds", sum)
		return nil
	}, time.Second)

	dsManager := datasources.NewDataSourceRuntimeManagerCommon(uuid.New(), "Teste", testLogger)
	dsManager.WithWorker(worker)

	rt.AddDataSourceManager(dsManager)

	rt.RunAll(ctx)

	<-time.After(5 * time.Second)

	rt.StopAll(ctx)

	<-time.After(1 * time.Second)

	if sum == 5 {
		t.Errorf("Falied sum, got %d", sum)
	}

	for _, manager := range rt.dataSources {
		if manager.Status() != datasources.Stopped {
			t.Errorf("expected to shutdown all data sources")
		}
	}

}
