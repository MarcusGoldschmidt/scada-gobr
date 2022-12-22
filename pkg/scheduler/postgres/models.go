package postgres

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/scheduler"
	"time"
)

type JobModel struct {
	Id            string          `gorm:"column:id;primary_key"`
	JobId         string          `gorm:"column:job_id"`
	TypeName      string          `gorm:"column:type_name"`
	Cron          string          `gorm:"column:cron"`
	NextExecution *time.Time      `gorm:"column:next_execution"`
	ScheduledJobs []*ScheduledJob `gorm:"foreignKey:JobId"`
}

func NewJobModelFromEntity(job *scheduler.JobEntity) *JobModel {
	return &JobModel{
		Id:            job.Id,
		JobId:         job.JobId,
		TypeName:      job.TypeName,
		Cron:          job.Cron,
		NextExecution: job.NextExecution,
	}
}

func (j *JobModel) ToEntity() *scheduler.JobEntity {
	return &scheduler.JobEntity{
		Id:            j.Id,
		JobId:         j.JobId,
		TypeName:      j.TypeName,
		Cron:          j.Cron,
		NextExecution: j.NextExecution,
	}
}

type ScheduledJob struct {
	Id         string     `gorm:"column:id;primary_key"`
	JobId      string     `gorm:"column:job_id"`
	TypeName   string     `gorm:"column:type_name"`
	At         time.Time  `gorm:"column:at"`
	ExecutedAt *time.Time `gorm:"column:executed_at"`
}

func (j ScheduledJob) ToEntity() *scheduler.ScheduledJob {
	return &scheduler.ScheduledJob{
		Id:         j.Id,
		JobId:      j.JobId,
		TypeName:   j.TypeName,
		At:         j.At,
		ExecutedAt: j.ExecutedAt,
	}
}

func NewScheduledJobFromModel(model *JobModel, executedAt *time.Time) *ScheduledJob {
	return &ScheduledJob{
		Id:         model.Id,
		JobId:      model.JobId,
		TypeName:   model.TypeName,
		At:         model.NextExecution.UTC(),
		ExecutedAt: executedAt,
	}
}
