package ran

import (
	"fmt"
)

type CommandRunner interface {
	RunCommand(command string, renv RuntimeEnvironment) error
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

func (cr StdCommandRunner) RunCommand(command string, renv RuntimeEnvironment) error {
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
			cr,
			supervisor,
			dispatcher,
			stack,
			cr.logger,
			renv,
		)
		if len(task.When) == 0 {
			initialRunners = append(initialRunners, tr)
		} else {
			dispatcher.Register(tr)
		}
	}

	for _, tr := range initialRunners {
		tr.Run()
	}

	resultErr := supervisor.Wait()
	for {
		cmd, ok := stack.Pop()
		if !ok {
			break
		}
		if err := cmd.Run(); err != nil {
			if resultErr == nil {
				resultErr = err
			} else {
				cr.logger.Error(err.Error())
			}
		}
	}
	return resultErr
}
