package scheduler

import (
	"errors"
	"time"
)

type ScheduledJob struct {
	Id         string
	JobId      string
	TypeName   string
	At         time.Time
	ExecutedAt *time.Time
}

type JobEntity struct {
	Id            string
	JobId         string
	TypeName      string
	Cron          string
	NextExecution *time.Time
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
