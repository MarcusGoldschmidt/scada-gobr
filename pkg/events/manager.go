package events

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"sync"
)

type HubClient interface {
	Execute(ctx context.Context, message any) error
}

type message struct {
	ctx     context.Context
	message any
}

type HubManager interface {
	SendMessage(ctx context.Context, topic string, data any)

	ShutDown(ctx context.Context)
	AddClient(ctx context.Context, topic string, client HubClient)
	RemoveClient(topic string, client HubClient)
	CreateTopic(ctx context.Context, topicName string)
}

type HubManagerImpl struct {
	mutex  sync.RWMutex
	topics map[string]*Hub
	logger logger.Logger
}

func NewHubManagerImpl(logger logger.Logger) *HubManagerImpl {
	return &HubManagerImpl{
		mutex:  sync.RWMutex{},
		topics: map[string]*Hub{},
		logger: logger,
	}
}

//ShutDown and remove all clients
func (hm *HubManagerImpl) ShutDown(ctx context.Context) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(hm.topics))

	for _, hub := range hm.topics {
		go func(hub *Hub) {
			hub.closed <- true
			wg.Done()
		}(hub)
	}
	wg.Wait()

	hm.topics = map[string]*Hub{}
}

func (hm *HubManagerImpl) SendMessage(ctx context.Context, topic string, data any) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	if hub, ok := hm.topics[topic]; ok {
		go func(hub *Hub) {
			hub.broadcast <- message{ctx, data}
		}(hub)
	}
}

func (hm *HubManagerImpl) AddClient(ctx context.Context, topic string, client HubClient) {
	hm.CreateTopic(ctx, topic)
	hm.topics[topic].register <- client
}

func (hm *HubManagerImpl) RemoveClient(topic string, client HubClient) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if hub, ok := hm.topics[topic]; ok {
		hub.unregister <- client
	}
}

func (hm *HubManagerImpl) CreateTopic(ctx context.Context, topicName string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	if _, ok := hm.topics[topicName]; !ok {
		hm.topics[topicName] = NewHub(hm.logger)
		go hm.topics[topicName].Run(ctx)
	}
}

type Hub struct {
	lock sync.RWMutex

	// Registered clients.
	clients map[HubClient]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan HubClient

	// Unregister requests from clients.
	unregister chan HubClient

	logger logger.Logger

	closed chan bool
}

func NewHub(logger logger.Logger, bufferSize ...int) *Hub {
	size := 256

	if len(bufferSize) > 0 {
		size = bufferSize[0]
	}

	return &Hub{
		lock:       sync.RWMutex{},
		broadcast:  make(chan message, size),
		clients:    make(map[HubClient]bool),
		register:   make(chan HubClient),
		unregister: make(chan HubClient),
		closed:     make(chan bool),
		logger:     logger,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-h.closed:
			return
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.lock.Lock()
			h.clients[client] = true
			h.lock.Unlock()
		case client := <-h.unregister:
			h.lock.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
			h.lock.Unlock()
		case message := <-h.broadcast:
			h.lock.RLock()
			for client := range h.clients {
				go func(client HubClient) {
					err := client.Execute(message.ctx, message.message)
					if err != nil {
						h.logger.Errorf("Error executing event client: %s", err)
					}
				}(client)
			}
			h.lock.RUnlock()
		}
	}
}
