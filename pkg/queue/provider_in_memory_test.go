package queue

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"testing"
)

func TestAck(t *testing.T) {
	ctx := context.Background()
	queue := NewProviderInMemory(providers.DefaultTimeProvider)

	_ = queue.Enqueue(ctx, "test", "test")

	msgs, _ := queue.Dequeue(ctx, "test", 1)

	msg := msgs[0]

	queue.Ack(ctx, "test", msg.Id())

	if queue.queueAck["test"+msg.Id()].ackTime == nil {
		t.Error("Message was not acked")
	}
}

func TestNack(t *testing.T) {
	ctx := context.Background()
	queue := NewProviderInMemory(providers.DefaultTimeProvider)

	_ = queue.Enqueue(ctx, "test", "test")

	msgs, _ := queue.Dequeue(ctx, "test", 1)

	msg := msgs[0]

	queue.Nack(ctx, "test", msg.Id(), errors.New("test"))

	if queue.queueAck["test"+msg.Id()].ackTime != nil {
		t.Error("Message was acked")
	}

	if len(queue.queueAck["test"+msg.Id()].errors) != 1 {
		t.Error("Message was not nacked")
	}
}
