package ran

type Dispatcher struct {
	logger    Logger
	receivers []EventReceiver
}

func NewDispatcher(logger Logger) *Dispatcher {
	return &Dispatcher{logger, nil}
}

func (d *Dispatcher) Receive(e Event) {
	d.logger.Debug("event: %q", e.Topic)
	for _, receiver := range d.receivers {
		receiver.Receive(e)
	}
}

func (d *Dispatcher) Register(tr *TaskRunner) {
	d.receivers = append(d.receivers, tr)
}
