package scheduler

import (
	"errors"
	"time"
)

type ScheduledJob struct {
	id         string
	jobId      string
	typeName   string
	at         time.Time
	executedAt *time.Time
}

func (s ScheduledJob) Id() string {
	return s.id
}

func (s ScheduledJob) JobId() string {
	return s.jobId
}

func (s ScheduledJob) TypeName() string {
	return s.typeName
}

// At returns the time when the job should be executed
func (s ScheduledJob) At() time.Time {
	return s.at
}

func (s ScheduledJob) ExecutedAt() *time.Time {
	return s.executedAt
}

type JobEntity struct {
	Id            string
	JobId         string
	TypeName      string
	Cron          string
	NextExecution *time.Time
}

func (j *JobEntity) name() {

}

func (j *JobEntity) validate() error {
	if j.Id == "" {
		return errors.New("id is required")
	}

	if j.JobId == "" {
		return errors.New("jobId is required")
	}

	if j.TypeName == "" {
		return errors.New("typeName is required")
	}

	if j.NextExecution == nil {
		if j.Cron == "" {
			return errors.New("cron is required")
		}
	}

	return nil
}
