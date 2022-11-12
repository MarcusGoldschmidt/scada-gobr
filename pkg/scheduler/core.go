package scheduler

import (
	"context"
	"time"
)

type Provider interface {
	// AddOrUpdateJob add or update a job, return true if the job was added and false if the job was updated
	AddOrUpdateJob(ctx context.Context, job *JobEntity) error

	ScheduleJobTime(ctx context.Context, date time.Time, jobId string) error
	RemoveScheduledJobs(ctx context.Context, jobId string) error

	GetJobsPeriod(ctx context.Context, initial, end time.Time) ([]*ScheduledJob, error)
	GetJobsPeriodById(ctx context.Context, jobId string, initial time.Time, end time.Time) ([]*ScheduledJob, error)

	RemoveJob(ctx context.Context, jobId string) error

	GetJobById(ctx context.Context, jobId string) (*JobEntity, error)
}
