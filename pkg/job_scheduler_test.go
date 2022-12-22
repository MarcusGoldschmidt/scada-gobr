package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/scheduler"
	"testing"
	"time"
)

func newSchedulerManager(t *testing.T) *SchedulerManager {
	scada := &Scadagobr{}

	queueManager, _ := createQueueManagerTest(t)
	queueManager.WithScada(scada)

	go queueManager.Start(context.Background())

	return NewSchedulerManager(scheduler.NewInMemoryProvider(), providers.DefaultTimeProvider, logger.NewTestLogger(t), queueManager)
}

func TestStartAndStop(t *testing.T) {
	manager := newSchedulerManager(t)

	go manager.Start(context.Background())

	time.Sleep(time.Millisecond * 100)

	err := manager.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

type dummyBackGroundJob struct {
}

func (d dummyBackGroundJob) Execute(ctx context.Context, scada *Scadagobr, data any) error {
	return nil
}

func TestStartScheduleAndTigger(t *testing.T) {
	manager := newSchedulerManager(t)

	go manager.Start(context.Background(), time.Second)

	time.Sleep(time.Millisecond * 100)

	err := manager.AddOrUpdateJob(context.Background(), "* * * * *", "teste", dummyBackGroundJob{})
	if err != nil {
		t.Fatal(err)
	}

	err = manager.TriggerJob(context.Background(), "teste")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
}

func TestScheduleJob(t *testing.T) {
	ctx := context.Background()
	timeProvider := providers.DefaultTimeProvider

	manager := newSchedulerManager(t)

	go manager.Start(ctx, time.Second)

	time.Sleep(time.Millisecond * 100)

	err := manager.AddOrUpdateJob(ctx, "* * * * *", "teste", dummyBackGroundJob{})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 1)

	period, err := manager.SchedulerProvider.GetJobsPeriod(ctx, timeProvider.GetCurrentTime(), timeProvider.GetCurrentTime().Add(time.Minute*60))
	if err != nil {
		t.Fatal(err)
	}

	if len(period) == 1 {
		t.Fatal("job not scheduled")
	}
}
