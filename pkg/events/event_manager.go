package events

import (
	"encoding/json"
	"io"
	"time"
)

type Event struct {
	time time.Time
	data []byte
}

type EventManager struct {
	io     io.ReadWriter
	events []*Event
}

func (ev *EventManager) AddEvent(data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	ev.events = append(ev.events, &Event{time.Now(), out})
}

func (ev *EventManager) GetEvents() []*Event {
	return ev.events
}
