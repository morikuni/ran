package ran

import (
	"context"
	"fmt"
)

type CommandRunner interface {
	RunCommand(ctx context.Context, command string) error
}

type StdCommandRunner struct {
	commands map[string]Command

	logger Logger
}

func NewStdCommandRunner(
	commands map[string]Command,
	logger Logger,
) StdCommandRunner {
	return StdCommandRunner{
		commands,
		logger,
	}
}

func (cr StdCommandRunner) RunCommand(ctx context.Context, command string, env RuntimeEnvironment) error {
	cmd, ok := cr.commands[command]
	if !ok {
		return fmt.Errorf("no such command: %s", command)
	}

	supervisor := NewSupervisor()
	dispatcher := NewDispatcher(cr.logger)
	stack := NewStack()

	var initialRunners []*TaskRunner
	for _, task := range cmd.Tasks {
		tr := NewTaskRunner(
			task,
			supervisor,
			dispatcher,
			stack,
			cr.logger,
			env,
		)
		if len(task.When) == 0 {
			initialRunners = append(initialRunners, tr)
		} else {
			dispatcher.Register(ctx, tr)
		}
	}

	for _, tr := range initialRunners {
		tr.Run(ctx)
	}

	if err := supervisor.Wait(); err != nil {
		return err
	}

	for {
		cmd, ok := stack.Pop()
		if !ok {
			break
		}
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
