package workflow

import (
	"context"
	"strings"
)

type Event struct {
	Topic   string
	Payload map[string]string
}

func NewEvent(topic string, payload map[string]string) Event {
	return Event{topic, payload}
}

type EventHandler interface {
	Receive(ctx context.Context, e Event) error
}
