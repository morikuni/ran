package ran

import (
	"context"
	"sync"
)

type WorkerStarter interface {
	Start(ctx context.Context, f func(ctx context.Context) error)
}

type Supervisor struct {
	wg sync.WaitGroup

	mu      sync.Mutex
	lastErr error
}

var _ interface {
	WorkerStarter
} = (*Supervisor)(nil)

func NewSupervisor() *Supervisor {
	return &Supervisor{}
}

func (s *Supervisor) Start(ctx context.Context, f func(ctx context.Context) error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := f(ctx)
		s.mu.Lock()
		s.lastErr = err
		s.mu.Unlock()
	}()
}

func (s *Supervisor) Wait() error {
	s.wg.Wait()
	return s.lastErr
}
