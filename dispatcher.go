package workflow

import (
	"context"
)

type Dispatcher struct {
	logger    Logger
	receivers []EventReceiver
}

func NewDispatcher(logger Logger) *Dispatcher {
	return &Dispatcher{logger, nil}
}

func (d *Dispatcher) Receive(ctx context.Context, e Event) error {
	d.logger.Info(e.Topic)
	for _, receiver := range d.receivers {
		if err := receiver.Receive(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dispatcher) Register(ctx context.Context, tr *TaskRunner) {
	d.receivers = append(d.receivers, tr)
}
