package runtime

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence/in_memory"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()

	log := logger.NewSimpleLogger("teste", os.Stdout)

	rt := NewRuntimeManager(log, in_memory.NewInMemoryPersistence(), events.NewHubManagerImpl(log), providers.DefaultTimeProvider)

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

	err := rt.RunAll(ctx)
	if err != nil {
		t.Errorf("Error running runtime: %s", err)
		return
	}

	<-time.After(5 * time.Second)

	rt.StopAll(ctx)

	<-time.After(1 * time.Second)

	if sum == 5 {
		t.Errorf("Falied sum, got %d", sum)
	}

	for _, manager := range rt.dataSources {
		if manager.dataSourceRuntimeManager.Status() != datasources.Stopped {
			t.Errorf("expected to shutdown all data sources")
		}
	}

}
