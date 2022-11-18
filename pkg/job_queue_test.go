package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/queue"
	"testing"
	"time"
)

func createQueueManagerTest(t *testing.T) (*QueueManager, queue.Provider) {
	provider := queue.NewProviderInMemory(providers.DefaultTimeProvider)

	return NewManager(provider, logger.NewTestLogger(t)), provider
}

func TestManagerStartAndStop(t *testing.T) {
	ctx := context.Background()
	manager, _ := createQueueManagerTest(t)

	go manager.Start(ctx)

	time.Sleep(time.Millisecond * 100)

	manager.Stop()
}

func TestManagerAddJobThenStop(t *testing.T) {
	ctx := context.Background()
	manager, _ := createQueueManagerTest(t)

	go manager.Start(ctx)

	manager.AddListener("test", "test", JobFunc(func(ctx context.Context, scada *Scadagobr, data any) error {
		return nil
	}))

	time.Sleep(time.Millisecond * 200)

	status, err := manager.GetJobsStatus("test")
	if err != nil {
		t.Error(err)
		return
	}

	if status[0].Status() != queue.Idle {
		t.Errorf("Expected status to be idle")
		return
	}

	manager.Stop()

	time.Sleep(time.Millisecond * 100)

	status, err = manager.GetJobsStatus("test")
	if err != nil {
		t.Error(err)
		return
	}

	if status[0].Status() != queue.Done {
		t.Errorf("Expected status to be done got " + status[0].Status().String())
		return
	}
}

func TestManagerAddJob(t *testing.T) {
	ctx := context.Background()
	manager, queueProvider := createQueueManagerTest(t)

	go manager.Start(ctx)

	checkValue := 0

	manager.AddListener("test", "test", JobFunc(func(ctx context.Context, scada *Scadagobr, data any) error {
		checkValue = data.(int)

		return nil
	}))

	err := queueProvider.Enqueue(ctx, "test", 1)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	if checkValue == 0 {
		t.Errorf("Expected change check value")
	}
}
