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

func (r *EventRecorder) GetValue(index int, path string) string {
	if len(r.Events) <= index {
		return ""
	}
	return r.Events[index].Payload[path]
}

func (r *EventRecorder) GetTopic(index int) string {
	if len(r.Events) <= index {
		return ""
	}
	return r.Events[index].Topic
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
