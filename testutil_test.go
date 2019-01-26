package ran_test

import (
	"context"

	"github.com/morikuni/ran"
)

type EventRecorder struct {
	Events []ran.Event
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{}
}

func (r *EventRecorder) Receive(ctx context.Context, e ran.Event) {
	r.Events = append(r.Events, e)
}

type SynchronousStarter struct {
	Error error
}

func NewSynchronousStarter() SynchronousStarter {
	return SynchronousStarter{}
}

func (s SynchronousStarter) Start(ctx context.Context, f func(ctx context.Context) error) {
	s.Error = f(ctx)
}

type CommandRecorder struct {
	Commands []string
}

func NewCommandRecorder() *CommandRecorder {
	return &CommandRecorder{}
}

func (r *CommandRecorder) RunCommand(ctx context.Context, command string, renv ran.RuntimeEnvironment) error {
	r.Commands = append(r.Commands, command)
	return nil
}
