package pkg

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/queue"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/scheduler"
	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"reflect"
	"sync"
	"time"
)

var schedulerTypeRegistry = map[string]reflect.Type{}

// erros
var (
	ErrJobNotFound         = errors.New("job not found")
	ErrJobAlreadyScheduled = errors.New("job already scheduled")
)

type SchedulerManager struct {
	SchedulerProvider scheduler.Provider

	lock         sync.Locker
	cancel       func()
	wg           *sync.WaitGroup
	timeProvider providers.TimeProvider
	logger       logger.Logger

	queueManager *QueueManager
}

func NewSchedulerManager(provider scheduler.Provider, timeProvider providers.TimeProvider, logger logger.Logger, queueManager *QueueManager) *SchedulerManager {
	return &SchedulerManager{
		SchedulerProvider: provider,
		timeProvider:      timeProvider,
		logger:            logger,
		queueManager:      queueManager,
		lock:              &sync.Mutex{},
		wg:                &sync.WaitGroup{},
	}
}

// AddOrUpdateJob adds a new job to the scheduler, the job must be a struct that implements the Job interface
func (m *SchedulerManager) AddOrUpdateJob(ctx context.Context, cron string, jobId string, job Job) error {
	cronParse, err := cronexpr.Parse(cron)

	if err != nil {
		return err
	}

	jobDatabase, err := m.SchedulerProvider.GetJobById(ctx, jobId)

	if err != nil {
		return err
	}

	jobType := reflect.TypeOf(job)

	entity := &scheduler.JobEntity{
		Id:       uuid.New().String(),
		Cron:     cron,
		JobId:    jobId,
		TypeName: jobType.String(),
	}

	schedulerTypeRegistry[entity.TypeName] = jobType

	// nothing changed
	if jobDatabase != nil && jobDatabase.Cron == cron && jobDatabase.TypeName == entity.TypeName {
		return nil
	}

	err = m.SchedulerProvider.AddOrUpdateJob(ctx, entity)
	if err != nil {
		return err
	}

	if jobDatabase != nil && jobDatabase.Cron != cron {
		err := m.SchedulerProvider.RemoveScheduledJobs(ctx, jobId)
		if err != nil {
			return err
		}
	}

	// new job or updated cron
	if jobDatabase == nil || jobDatabase.Cron != cron {
		times := cronParse.NextN(m.timeProvider.GetCurrentTime(), 30)

		for _, t := range times {
			err := m.SchedulerProvider.ScheduleJobTime(ctx, t, jobId)
			if err != nil {
				return err
			}
		}
	}

	m.queueManager.UpdateListener("scheduler:"+entity.TypeName, jobId, createCallBackSchedule(entity))

	return nil
}

func (m *SchedulerManager) TriggerJob(ctx context.Context, jobId string) error {
	job, err := m.SchedulerProvider.GetJobById(ctx, jobId)
	if err != nil {
		return err
	}

	err = m.queueProvider().Enqueue(ctx, "scheduler:"+job.TypeName, struct{}{})
	if err != nil {
		return err
	}

	return nil
}

func createCallBackSchedule(entity *scheduler.JobEntity) Job {
	return JobFunc(func(ctx context.Context, scada *Scadagobr, data any) error {
		return reflect.New(schedulerTypeRegistry[entity.TypeName]).Interface().(Job).Execute(ctx, scada, data)
	})
}

func (m *SchedulerManager) RemoveJob(ctx context.Context, jobId string) error {

	err := m.SchedulerProvider.RemoveScheduledJobs(ctx, jobId)
	if err != nil {
		return err
	}

	err = m.SchedulerProvider.RemoveJob(ctx, jobId)
	if err != nil {
		return err
	}

	return nil
}

func (m *SchedulerManager) GetNextJobs(ctx context.Context, jobId string, duration time.Duration) ([]*scheduler.ScheduledJob, error) {
	now := m.timeProvider.GetCurrentTime()

	nextJobs, err := m.SchedulerProvider.GetJobsPeriodById(ctx, jobId, now, now.Add(duration))
	if err != nil {
		return nil, err
	}

	return nextJobs, err
}

func (m *SchedulerManager) queueProvider() queue.Provider {
	return m.queueManager.QueueProvider
}

func (m *SchedulerManager) Stop() error {
	m.cancel()

	return nil
}

func (m *SchedulerManager) ScheduleNextExecution(ctx context.Context, jobId string) error {
	entity, err := m.SchedulerProvider.GetJobById(ctx, jobId)

	if err != nil {
		return err
	}

	var nextTime time.Time

	if entity.NextExecution != nil {
		nextTime = *entity.NextExecution
	} else if entity.Cron != "" {
		cronParse, err := cronexpr.Parse(entity.Cron)

		if err != nil {
			return err
		}

		nextTime = cronParse.Next(m.timeProvider.GetCurrentTime())
	} else {
		return errors.New("no cron or next execution time found for job: " + entity.JobId)
	}

	jobExist, err := m.SchedulerProvider.GetJobsTimeAndId(ctx, entity.JobId, nextTime)
	if err != nil {
		return err
	}

	if jobExist != nil {
		return ErrJobAlreadyScheduled
	}

	return m.SchedulerProvider.ScheduleJobTime(ctx, nextTime, entity.JobId)
}

func (m *SchedulerManager) Start(ctx context.Context, verificationTimes ...time.Duration) {
	verificationTime := time.Minute

	if len(verificationTimes) > 0 {
		verificationTime = verificationTimes[0]
	}

	ticker := time.Tick(verificationTime)

	ctx, m.cancel = context.WithCancel(ctx)

	for {
	SELECT:
		select {
		case <-ctx.Done():
			m.wg.Wait()
			return
		case <-ticker:
			m.lock.Lock()

			startTime := m.timeProvider.GetCurrentTime().Add(-verificationTime)

			jobs, err := m.SchedulerProvider.GetJobsPeriod(ctx, startTime, startTime.Add(verificationTime*2))
			if err != nil {
				m.logger.Errorf("Failed to get next jobs: %s", err)
				m.lock.Unlock()
				goto SELECT
			}

			m.wg.Add(len(jobs))

			for _, job := range jobs {
				go func(job *scheduler.ScheduledJob) {
					defer m.wg.Done()

					now := m.timeProvider.GetCurrentTime()

					// Sleep until the job should be executed
					time.Sleep(job.At.Sub(now))

					err := m.queueProvider().Enqueue(ctx, "scheduler:"+job.TypeName, struct{}{})
					if err != nil {
						m.logger.Errorf("Failed to enqueue job: %s", err)
					}

					err = m.ScheduleNextExecution(ctx, job.JobId)
					if err != nil && err != ErrJobAlreadyScheduled {
						m.logger.Errorf("Failed to schedule next execution: %s", err)
					}
				}(job)
			}

			m.wg.Wait()

			m.lock.Unlock()
		}
	}
}
