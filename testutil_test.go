package ran_test

import (
	"github.com/morikuni/ran"
)

type EventRecorder struct {
	Events []ran.Event
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{}
}

func (r *EventRecorder) Receive(e ran.Event) {
	r.Events = append(r.Events, e)
}

type SynchronousStarter struct {
	Error error
}

func NewSynchronousStarter() SynchronousStarter {
	return SynchronousStarter{}
}

func (s SynchronousStarter) Start(f func() error) {
	s.Error = f()
}

type CommandRecorder struct {
	Commands []string
}

func NewCommandRecorder() *CommandRecorder {
	return &CommandRecorder{}
}

func (r *CommandRecorder) RunCommand(command string, renv ran.RuntimeEnvironment) error {
	r.Commands = append(r.Commands, command)
	return nil
}
