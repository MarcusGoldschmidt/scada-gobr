package postgres

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/test"
	"testing"
)

func setupProvider(container *test.SqlContainer, t *testing.T) *SqlPostgresJobQueue {
	provider := NewSqlPostgresJobQueue(container.GormDb(), providers.DefaultTimeProvider)

	err := provider.Setup()
	if err != nil {
		t.Fatalf("Error while migrating: %s", err)
	}
	return provider
}

func TestEnqueueAck(t *testing.T) {
	container := test.SetupPostgresContainer(t)
	defer container.Close()

	provider := setupProvider(container, t)

	ctx := context.Background()

	err := provider.Enqueue(ctx, "test", 1)
	if err != nil {
		t.Fatalf("Error while enqueueing: %s", err)
	}

	err = provider.Enqueue(ctx, "test", 2)
	if err != nil {
		t.Fatalf("Error while enqueueing: %s", err)
	}

	dequeue, err := provider.Dequeue(ctx, "test", 10)
	if err != nil {
		t.Fatalf("Error while dequeueing: %s", err)
	}

	if len(dequeue) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(dequeue))
	}

	for _, message := range dequeue {
		provider.Ack(ctx, "test", message.Id())
	}

	dequeue, err = provider.Dequeue(ctx, "test", 10)
	if err != nil {
		t.Fatalf("Error while dequeueing: %s", err)
	}

	if len(dequeue) != 0 {
		t.Fatalf("Expected 0 messages, got %d", len(dequeue))
	}
}

func TestEnqueueNack(t *testing.T) {
	container := test.SetupPostgresContainer(t)
	defer container.Close()

	provider := setupProvider(container, t)

	ctx := context.Background()

	err := provider.Enqueue(ctx, "test", 1)
	if err != nil {
		t.Fatalf("Error while enqueueing: %s", err)
	}

	dequeue, err := provider.Dequeue(ctx, "test", 10)
	if err != nil {
		t.Fatalf("Error while dequeueing: %s", err)
	}

	if len(dequeue) != 1 {
		t.Fatalf("Expected 1 messages, got %d", len(dequeue))
	}

	for _, message := range dequeue {
		provider.Nack(ctx, "test", message.Id(), errors.New("test error"))
	}

	dequeue, err = provider.Dequeue(ctx, "test", 10)
	if err != nil {
		t.Fatalf("Error while dequeueing: %s", err)
	}

	if len(dequeue) != 1 {
		t.Fatalf("Expected 1 messages, got %d", len(dequeue))
	}
}
