package scheduler

import (
	"context"
	"time"
)

type Provider interface {
	AddOrUpdateJob(ctx context.Context, job *JobEntity) error
	RemoveJob(ctx context.Context, jobId string) error

	ScheduleJobTime(ctx context.Context, date time.Time, jobId string) error
	RemoveScheduledJobs(ctx context.Context, jobId string) error

	GetJobsPeriod(ctx context.Context, initial, end time.Time) ([]*ScheduledJob, error)
	GetJobsTimeAndId(ctx context.Context, jobId string, time time.Time) (*ScheduledJob, error)
	GetJobsPeriodById(ctx context.Context, jobId string, initial time.Time, end time.Time) ([]*ScheduledJob, error)
	GetJobById(ctx context.Context, jobId string) (*JobEntity, error)
}
