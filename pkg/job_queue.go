package pkg

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/queue"
	"sync"
	"time"
)

type jobAggregate struct {
	job    *jobInternal
	cancel func()
}

type queueAggregate struct {
	lock        sync.Locker
	sendMessage chan message
	jobs        []jobAggregate
}

type JobStatusResponse struct {
	status queue.JobStatus
	jobId  string
}

func (j JobStatusResponse) Status() queue.JobStatus {
	return j.status
}

func (j JobStatusResponse) JobId() string {
	return j.jobId
}

type QueueManager struct {
	scada         *Scadagobr
	lock          sync.Locker
	queues        map[string]*queueAggregate
	logger        logger.Logger
	queueProvider queue.Provider

	cancel func()
}

func NewManager(jobQueue queue.Provider, logger logger.Logger) *QueueManager {
	return &QueueManager{
		lock:          &sync.Mutex{},
		queues:        map[string]*queueAggregate{},
		queueProvider: jobQueue,
		logger:        logger,
	}
}

func (m *QueueManager) WithScada(scada *Scadagobr) {
	m.scada = scada
}

func (m *QueueManager) Start(ctx context.Context) {
	ticker := time.Tick(time.Millisecond * 200)

	ctx, m.cancel = context.WithCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			m.lock.Lock()

			for queueName, _ := range m.queues {
				msgs, err := m.queueProvider.Dequeue(ctx, queueName, 16)
				if err != nil {
					return
				}

				for _, msg := range msgs {
					internalMsg := message{
						ctx:      msg.Ctx(),
						id:       msg.Id(),
						data:     msg.Data(),
						response: make(chan error),
					}

					m.queues[queueName].sendMessage <- internalMsg

					go func(queue string, msg message) {
						err := <-msg.response
						if err != nil {
							m.logger.Errorf("Queue %s MessageId %s failed: %s", queue, msg.id, err)
							m.queueProvider.Nack(ctx, queueName, msg.id, err)
						} else {
							m.queueProvider.Ack(ctx, queueName, msg.id)
						}
					}(queueName, internalMsg)
				}
			}

			m.lock.Unlock()
		}
	}
}

func (m *QueueManager) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(m.queues))

	for _, queueAgg := range m.queues {
		go func(queueAgg *queueAggregate) {
			defer wg.Done()

			queueAgg.lock.Lock()

			for _, jobAgg := range queueAgg.jobs {
				jobAgg.cancel()
			}

			queueAgg.lock.Unlock()
		}(queueAgg)
	}

	wg.Wait()
}

func (m *QueueManager) UpdateListener(queueName string, jobId string, job Job) {
	m.RemoveListener(queueName, jobId)
	m.AddListener(queueName, jobId, job)
}

func (m *QueueManager) RemoveListener(queueName string, jobId string) {
	queueAgg := m.getOrCreateQueue(queueName)

	queueAgg.lock.Lock()
	defer queueAgg.lock.Unlock()

	jobs := queueAgg.jobs

	for i, jobAgg := range jobs {
		if jobAgg.job.id == jobId {
			jobAgg.cancel()

			queueAgg.jobs = append(queueAgg.jobs[:i], queueAgg.jobs[i+1:]...)
			break
		}
	}
}

func (m *QueueManager) AddListener(queueName string, jobId string, job Job) {
	queueAgg := m.getOrCreateQueue(queueName)

	queueAgg.lock.Lock()
	defer queueAgg.lock.Unlock()

	jobInternal := jobInternal{
		id:     jobId,
		logger: m.logger,
		queue:  queueName,
		job:    job,
		status: queue.Idle,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go jobInternal.Run(ctx, m.scada, queueAgg.sendMessage)

	queueAgg.jobs = append(queueAgg.jobs, jobAggregate{
		job:    &jobInternal,
		cancel: cancel,
	})
}

func (m *QueueManager) GetJobsStatus(queueName string) ([]*JobStatusResponse, error) {
	response := make([]*JobStatusResponse, 0)

	queueAgg, ok := m.queues[queueName]

	if !ok {
		return nil, errors.New("Queue not found " + queueName)
	}

	queueAgg.lock.Lock()
	defer queueAgg.lock.Unlock()

	for _, agg := range queueAgg.jobs {
		response = append(response, &JobStatusResponse{
			jobId:  agg.job.id,
			status: agg.job.status,
		})
	}

	return response, nil
}

func (m *QueueManager) getOrCreateQueue(queue string) *queueAggregate {
	m.lock.Lock()
	defer m.lock.Unlock()

	queueAgg, ok := m.queues[queue]
	if ok {
		return queueAgg
	}

	queueAgg = &queueAggregate{
		lock:        &sync.Mutex{},
		sendMessage: make(chan message, 64),
		jobs:        make([]jobAggregate, 0),
	}

	m.queues[queue] = queueAgg

	return queueAgg
}
