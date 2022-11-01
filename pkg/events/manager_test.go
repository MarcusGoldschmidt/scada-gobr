package events

import (
	"context"
	customLogger "github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"testing"
	"time"
)

type testClient struct {
	cb func(ctx context.Context, message any) error
}

func newTestClient(cb func(ctx context.Context, message any) error) *testClient {
	return &testClient{cb: cb}
}

func (t *testClient) Execute(ctx context.Context, message any) error {
	return t.cb(ctx, message)
}

func TestShouldAddClient(t *testing.T) {
	manager := NewHubManagerImpl(customLogger.NewTestLogger(t))

	ctx := context.Background()

	resultChan := make(chan int)

	testClient := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	manager.AddClient(ctx, "test", testClient)
	manager.SendMessage(ctx, "test", 1)

	select {
	case <-time.After(time.Second):
		t.Error("timeout")
	case <-resultChan:
		return
	}
}

func TestShouldAddAndRemoveClient(t *testing.T) {
	manager := NewHubManagerImpl(customLogger.NewTestLogger(t))

	ctx := context.Background()

	resultChan := make(chan int)

	testClient := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	manager.AddClient(ctx, "test", testClient)
	manager.RemoveClient("test", testClient)
	manager.SendMessage(ctx, "test", 1)

	select {
	case <-time.After(time.Second):
		return
	case <-resultChan:
		t.Error("should not receive message")
	}
}

func TestShouldBroadcastToClients(t *testing.T) {
	manager := NewHubManagerImpl(customLogger.NewTestLogger(t))

	ctx := context.Background()

	resultChan := make(chan int)

	testClient1 := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	testClient2 := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	manager.AddClient(ctx, "test", testClient1)
	manager.AddClient(ctx, "test", testClient2)
	manager.SendMessage(ctx, "test", 1)

	for i := 0; i < 2; i++ {
		select {
		case <-time.After(time.Second):
			t.Error("timeout", i)
		case <-resultChan:
			continue
		}
	}
}

func TestShouldBroadcastToSameClient(t *testing.T) {
	manager := NewHubManagerImpl(customLogger.NewTestLogger(t))

	ctx := context.Background()

	resultChan := make(chan int)

	testClient := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	manager.AddClient(ctx, "test1", testClient)
	manager.AddClient(ctx, "test2", testClient)
	manager.SendMessage(ctx, "test1", 1)
	manager.SendMessage(ctx, "test2", 1)

	for i := 0; i < 2; i++ {
		select {
		case <-time.After(time.Second):
			t.Error("timeout", i)
		case <-resultChan:
			continue
		}
	}
}

func TestShouldShutdown(t *testing.T) {
	manager := NewHubManagerImpl(customLogger.NewTestLogger(t))

	ctx := context.Background()

	resultChan := make(chan int)

	testClient := newTestClient(func(ctx context.Context, message any) error {
		resultChan <- 1
		return nil
	})

	manager.AddClient(ctx, "test", testClient)
	manager.AddClient(ctx, "test2", testClient)
	manager.ShutDown(ctx)
	manager.SendMessage(ctx, "test", 1)

	select {
	case <-time.After(time.Second):
		return
	case <-resultChan:
		t.Error("should not receive message")
	}
}
