package queue

import (
	"context"
)

//go:generate stringer -type=JobStatus -output=job_status_string.go
type JobStatus int8

const (
	Idle JobStatus = iota
	Running
	Done
)

type Message interface {
	Id() string
	Ctx() context.Context
	Data() any
}

type Provider interface {
	Enqueue(ctx context.Context, queue string, data any) error
	Dequeue(ctx context.Context, queue string, length uint) ([]Message, error)

	Ack(ctx context.Context, queue string, messageId string)
	Nack(ctx context.Context, queue string, messageId string, err error)
}
