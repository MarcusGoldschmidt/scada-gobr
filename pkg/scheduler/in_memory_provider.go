package scheduler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type InMemoryProvider struct {
	scheduledJobs map[string][]*ScheduledJob
	jobs          map[string]*JobEntity
}

func (i *InMemoryProvider) GetJobsTimeAndId(ctx context.Context, jobId string, time time.Time) (*ScheduledJob, error) {
	scheduledJobs, ok := i.scheduledJobs[jobId]

	if !ok {
		return nil, nil
	}

	for _, scheduledJob := range scheduledJobs {
		if scheduledJob.at.Equal(time) {
			return scheduledJob, nil
		}
	}

	return nil, nil
}

func NewInMemoryProvider() *InMemoryProvider {
	return &InMemoryProvider{
		scheduledJobs: map[string][]*ScheduledJob{},
		jobs:          map[string]*JobEntity{},
	}
}

func (i *InMemoryProvider) AddOrUpdateJob(ctx context.Context, job *JobEntity) error {
	i.jobs[job.JobId] = job

	return nil
}

func (i *InMemoryProvider) ScheduleJobTime(ctx context.Context, date time.Time, jobId string) error {
	job, ok := i.jobs[jobId]

	if !ok {
		return errors.New("job not found")
	}

	i.scheduledJobs[jobId] = append(i.scheduledJobs[jobId], &ScheduledJob{
		id:         uuid.New().String(),
		jobId:      jobId,
		at:         date,
		typeName:   job.TypeName,
		executedAt: nil,
	})

	return nil
}

func (i *InMemoryProvider) RemoveScheduledJobs(ctx context.Context, jobId string) error {
	i.scheduledJobs[jobId] = []*ScheduledJob{}

	return nil
}

func (i *InMemoryProvider) GetJobsPeriod(ctx context.Context, initial, end time.Time) ([]*ScheduledJob, error) {
	jobs := make([]*ScheduledJob, 0)

	for _, scheduledJobs := range i.scheduledJobs {
		for _, scheduledJob := range scheduledJobs {
			if scheduledJob.at.Before(end) && scheduledJob.at.After(initial) {
				jobs = append(jobs, scheduledJob)
			}
		}
	}

	return jobs, nil
}

func (i *InMemoryProvider) GetJobsPeriodById(ctx context.Context, jobId string, initial time.Time, end time.Time) ([]*ScheduledJob, error) {
	jobs := make([]*ScheduledJob, 0)

	scheduledJobs, ok := i.scheduledJobs[jobId]

	if !ok {
		return nil, errors.New("job not found")
	}

	for _, scheduledJob := range scheduledJobs {
		if scheduledJob.at.Before(end) && scheduledJob.at.After(initial) {
			jobs = append(jobs, scheduledJob)
		}
	}

	return jobs, nil
}

func (i *InMemoryProvider) RemoveJob(ctx context.Context, jobId string) error {
	delete(i.jobs, jobId)

	return nil
}

func (i *InMemoryProvider) GetJobById(ctx context.Context, jobId string) (*JobEntity, error) {
	job, _ := i.jobs[jobId]

	return job, nil
}
