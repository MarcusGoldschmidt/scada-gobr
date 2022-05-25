package runtime

import (
	"context"
	"os"
	"scadagobr/pkg/datasources"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/persistence"
	"testing"
	"time"
)

func TestSimpleRuntime(t *testing.T) {
	ctx := context.Background()

	rt := NewRuntimeManager(logger.NewSimpleLogger("teste", os.Stdout), &persistence.InMemoryPersistence{})

	sum := 0

	dataSource := datasources.NewCallbackDataSource(func() error {
		sum = sum + 1
		t.Logf("Pass %d seconds", sum)
		return nil
	}, true)

	rt.AddDataSource(dataSource)

	rt.Run(ctx, dataSource.Id())

	<-time.After(1 * time.Second)

	rt.StopAll()

	if len(rt.dataSourcesRuntimeControl) > 0 {
		t.Errorf("Expected to shutdown all data sources")
	}
}

func TestSimpleRuntimeSum(t *testing.T) {
	rt := NewRuntimeManager(logger.NewSimpleLogger("teste", os.Stdout), &persistence.InMemoryPersistence{})

	data := make(chan int)

	dataSource := datasources.NewCallbackDataSource(func() error {
		data <- 1
		return nil
	}, true)

	rt.AddDataSource(dataSource)

	ctx := context.Background()

	rt.Run(ctx, dataSource.Id())

	sum := 0

	go func() {
		<-time.After(5 * time.Second)
		rt.StopAll()
		close(data)
	}()

	for range data {
		sum = sum + 1
	}

	if sum == 5 {
		t.Errorf("Falied sum, got %d", sum)
	}

	if len(rt.dataSourcesRuntimeControl) != 0 {
		t.Errorf("Do not stop all datasources, len: %d", len(rt.dataSourcesRuntimeControl))
	}
}
