package workflow

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Supervisor struct {
	eg  *errgroup.Group
	ctx context.Context
}

func NewSupervisor(ctx context.Context) *Supervisor {
	eg, ctx := errgroup.WithContext(ctx)
	return &Supervisor{eg, ctx}
}

func (s *Supervisor) Start(f func(context.Context) error) {
	s.eg.Go(func() error {
		return f(s.ctx)
	})
}

func (s *Supervisor) Wait() error {
	return s.eg.Wait()
}
