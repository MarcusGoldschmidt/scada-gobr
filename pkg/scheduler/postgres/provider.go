package postgres

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/scheduler"
	"gorm.io/gorm"
	"time"
)

type SchedulerPostgresProvider struct {
	db *gorm.DB
}

func NewSchedulerPostgresProvider(db *gorm.DB) *SchedulerPostgresProvider {
	return &SchedulerPostgresProvider{db: db}
}

func (p *SchedulerPostgresProvider) Setup() error {
	err := p.db.Migrator().AutoMigrate(&JobModel{}, &ScheduledJob{})
	if err != nil {
		return err
	}

	return nil
}

func (p *SchedulerPostgresProvider) AddOrUpdateJob(ctx context.Context, job *scheduler.JobEntity) error {
	db := p.db.WithContext(ctx)

	model := NewJobModelFromEntity(job)

	var result *JobModel

	db.Model(JobModel{JobId: job.JobId}).First(&result)

	if result == nil {
		return db.Create(model).Error
	} else {
		return db.Model(result).Updates(model).Error
	}
}

func (p *SchedulerPostgresProvider) RemoveJob(ctx context.Context, jobId string) error {
	db := p.db.WithContext(ctx)

	return db.Delete(JobModel{JobId: jobId}).Error
}

func (p *SchedulerPostgresProvider) ScheduleJobTime(ctx context.Context, date time.Time, jobId string) error {
	db := p.db.WithContext(ctx)

	var job *JobModel
	db.Model(JobModel{JobId: jobId}).First(&job)

	if job == nil {
		return errors.New("job not found")
	}

	scheduledJob := NewScheduledJobFromModel(job, &date)

	return db.Save(scheduledJob).Error
}

func (p *SchedulerPostgresProvider) RemoveScheduledJobs(ctx context.Context, jobId string) error {
	db := p.db.WithContext(ctx)

	return db.Where("job_id = ?", jobId).Delete(ScheduledJob{}).Error
}

func (p *SchedulerPostgresProvider) GetJobsPeriod(ctx context.Context, initial, end time.Time) ([]*scheduler.ScheduledJob, error) {
	db := p.db.WithContext(ctx)

	var scheduledJobs []*ScheduledJob

	err := db.Where("at BETWEEN ? AND ?", initial, end).Find(&scheduledJobs).Error

	result := make([]*scheduler.ScheduledJob, len(scheduledJobs))
	for i, scheduledJob := range scheduledJobs {
		result[i] = scheduledJob.ToEntity()
	}

	return result, err
}

func (p *SchedulerPostgresProvider) GetJobsTimeAndId(ctx context.Context, jobId string, time time.Time) (*scheduler.ScheduledJob, error) {
	db := p.db.WithContext(ctx)

	var scheduledJob *ScheduledJob
	err := db.Where("job_id = ? AND at = ?", jobId, time).First(&scheduledJob).Error
	if err != nil {
		return nil, err
	}

	return scheduledJob.ToEntity(), nil
}

func (p *SchedulerPostgresProvider) GetJobsPeriodById(ctx context.Context, jobId string, initial time.Time, end time.Time) ([]*scheduler.ScheduledJob, error) {
	db := p.db.WithContext(ctx)

	var scheduledJobs []*ScheduledJob
	err := db.Where("job_id = ? AND at BETWEEN ? AND ?", jobId, initial, end).Find(&scheduledJobs).Error
	if err != nil {
		return nil, err
	}

	result := make([]*scheduler.ScheduledJob, len(scheduledJobs))
	for i, scheduledJob := range scheduledJobs {
		result[i] = scheduledJob.ToEntity()
	}

	return result, nil
}

func (p *SchedulerPostgresProvider) GetJobById(ctx context.Context, jobId string) (*scheduler.JobEntity, error) {
	db := p.db.WithContext(ctx)

	var job *JobModel
	err := db.Where("job_id = ?", jobId).First(&job).Error
	if err != nil {
		return nil, err
	}

	return job.ToEntity(), nil
}
