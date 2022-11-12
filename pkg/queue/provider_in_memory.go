package queue

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/google/uuid"
	"time"
)

type inMemoryMessage struct {
	id   string
	ctx  context.Context
	data any
}

func newInMemoryMessage(ctx context.Context, data any) *inMemoryMessage {
	return &inMemoryMessage{id: uuid.New().String(), ctx: ctx, data: data}
}

func (i *inMemoryMessage) Id() string {
	return i.id
}

func (i *inMemoryMessage) Ctx() context.Context {
	return i.ctx
}

func (i *inMemoryMessage) Data() any {
	return i.data
}

type nackError struct {
	err  error
	time time.Time
}

type queueAck struct {
	ackTime *time.Time
	errors  []nackError
}

type ProviderInMemory struct {
	queues map[string]chan *inMemoryMessage
	// Queue name + message id
	queueAck     map[string]*queueAck
	timeProvider providers.TimeProvider
}

func NewProviderInMemory(timeProvider providers.TimeProvider) *ProviderInMemory {
	return &ProviderInMemory{
		queues:       make(map[string]chan *inMemoryMessage),
		queueAck:     make(map[string]*queueAck),
		timeProvider: timeProvider,
	}
}

func (i *ProviderInMemory) getQueue(queueName string) chan *inMemoryMessage {
	queueChan, ok := i.queues[queueName]

	if !ok {
		queueChan = make(chan *inMemoryMessage, 12)
		i.queues[queueName] = queueChan
	}

	return queueChan
}

func (i *ProviderInMemory) Enqueue(ctx context.Context, queue string, data any) error {
	queueChan := i.getQueue(queue)

	go func() {
		msg := newInMemoryMessage(ctx, data)
		i.queueAck[queue+msg.id] = &queueAck{
			ackTime: nil,
			errors:  make([]nackError, 0),
		}
		queueChan <- msg
	}()

	return nil
}

func (i *ProviderInMemory) Dequeue(ctx context.Context, queueName string, length uint) ([]Message, error) {
	queueChan := i.getQueue(queueName)

	messages := make([]Message, 0)

	for {
		if length <= 0 {
			break
		}

		select {
		case <-time.After(time.Millisecond * 10):
			return messages, nil
		case msg := <-queueChan:
			messages = append(messages, msg)
		}

		length--
	}

	return messages, nil
}

func (i *ProviderInMemory) Ack(ctx context.Context, queue string, messageId string) {
	ack, ok := i.queueAck[queue+messageId]

	if !ok {
		return
	}

	now := i.timeProvider.GetCurrentTime()

	ack.ackTime = &now
}

func (i *ProviderInMemory) Nack(ctx context.Context, queue string, messageId string, err error) {
	ack, ok := i.queueAck[queue+messageId]

	if !ok {
		return
	}

	now := i.timeProvider.GetCurrentTime()

	ack.errors = append(ack.errors, nackError{
		err:  err,
		time: now,
	})
}
