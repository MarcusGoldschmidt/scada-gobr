package events

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"sync"
)

const (
	DataSeriesInserter = "dataseriesinserter:"
)

type HubClient interface {
	Execute(ctx context.Context, message any) error
}

type HubManager interface {
	Start(ctx context.Context)
	SendMessage(topic string, data any)
	AddClient(topic string, client HubClient)
	RemoveClient(topic string, client HubClient)
	CreateTopic(topicName string)
}

type HubManagerImpl struct {
	topics map[string]*Hub
	logger logger.Logger
}

func NewHubManagerImpl(logger logger.Logger) *HubManagerImpl {
	return &HubManagerImpl{
		topics: map[string]*Hub{},
		logger: logger,
	}
}

// Start starts the hub and can be stopped by cancel the contex
func (hm *HubManagerImpl) Start(ctx context.Context) {
	for topic, hub := range hm.topics {
		hm.logger.Infof("Stating hub for topic %s", topic)
		go hub.Run(ctx)
		hm.logger.Infof("Started hub for topic %s", topic)
	}
}

func (hm *HubManagerImpl) SendMessage(topic string, data any) {
	if hub, ok := hm.topics[topic]; ok {
		hub := hub
		go func() {
			hub.broadcast <- data
		}()
	}
}

func (hm *HubManagerImpl) AddClient(topic string, client HubClient) {
	hm.CreateTopic(topic)
	hm.topics[topic].register <- client
}

func (hm *HubManagerImpl) RemoveClient(topic string, client HubClient) {
	hm.CreateTopic(topic)
	hm.topics[topic].unregister <- client
}

func (hm *HubManagerImpl) CreateTopic(topicName string) {
	if _, ok := hm.topics[topicName]; !ok {
		hm.topics[topicName] = NewHub(hm.logger)
	}
}

type Hub struct {
	lock sync.RWMutex

	// Registered clients.
	clients map[HubClient]bool

	// Inbound messages from the clients.
	broadcast chan any

	// Register requests from the clients.
	register chan HubClient

	// Unregister requests from clients.
	unregister chan HubClient

	logger logger.Logger
}

func NewHub(logger logger.Logger, bufferSize ...int) *Hub {
	size := 16

	if len(bufferSize) > 0 {
		size = bufferSize[0]
	}

	return &Hub{
		lock:       sync.RWMutex{},
		broadcast:  make(chan any, size),
		clients:    map[HubClient]bool{},
		register:   make(chan HubClient),
		unregister: make(chan HubClient),
		logger:     logger,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
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
				client := client
				go func() {
					err := client.Execute(context.Background(), message)
					if err != nil {
						h.logger.Errorf("Error executing client: %s", err)
					}
				}()
			}
			h.lock.RUnlock()
		}
	}
}
