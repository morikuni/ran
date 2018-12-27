package workflow_test

import (
	"context"

	"github.com/morikuni/workflow"
)

type EventRecorder struct {
	Events []workflow.Event
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{}
}

func (r *EventRecorder) Receive(ctx context.Context, e workflow.Event) error {
	r.Events = append(r.Events, e)
	return nil
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
}

func NewSynchronousStarter() SynchronousStarter {
	return SynchronousStarter{}
}

func (s SynchronousStarter) Start(ctx context.Context, f func(ctx context.Context) error) {
	if err := f(ctx); err != nil {
		panic(err)
	}
}
