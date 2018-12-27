package workflow

import (
	"context"
)

type Dispatcher struct {
	logger Logger
}

func NewDispatcher(logger Logger) *Dispatcher {
	return &Dispatcher{logger}
}

func (d *Dispatcher) Receive(ctx context.Context, e Event) error {
	d.logger.Info("event: %#v", e.Topic)
	return nil
}
