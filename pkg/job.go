package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/queue"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/util"
	"go.opentelemetry.io/otel/attribute"
)

type Job interface {
	Execute(ctx context.Context, scada *Scadagobr, data any) error
}

type JobFunc func(ctx context.Context, scada *Scadagobr, data any) error

func (f JobFunc) Execute(ctx context.Context, scada *Scadagobr, data any) error {
	return f(ctx, scada, data)
}

type message struct {
	ctx  context.Context
	id   string
	data any

	response chan error
}

type jobInternal struct {
	id     string
	queue  string
	job    Job
	logger logger.Logger
	status queue.JobStatus
}

func (j *jobInternal) Run(ctx context.Context, scada *Scadagobr, data chan message) {
	for {
		j.status = queue.Idle

		select {
		case <-ctx.Done():
			j.status = queue.Done
			return
		case msg := <-data:
			j.status = queue.Running

			ctx, span := util.Tracer.Start(msg.ctx, "Job.Run")
			span.SetAttributes(
				attribute.KeyValue{
					Key:   "job.id",
					Value: attribute.StringValue(j.id),
				},
				attribute.KeyValue{
					Key:   "job.queue",
					Value: attribute.StringValue(j.queue),
				},
			)

			err := j.job.Execute(ctx, scada, msg.data)
			msg.response <- err
			span.End()
		}
	}
}
