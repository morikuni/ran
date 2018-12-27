package workflow

import (
	"context"
	"flag"
	"fmt"
	"os"
)

type App struct {
	logger Logger
}

func NewApp(logger Logger) App {
	return App{logger}
}

func (app App) Run(ctx context.Context, args []string, signal <-chan os.Signal) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	file := fs.String("f", "workflow.yaml", "file")
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	def, err := LoadDefinition(*file)
	if err != nil {
		return err
	}

	target := fs.Arg(0)
	command, ok := def.Commands[target]
	if !ok {
		return fmt.Errorf("no such workflow: %s", target)
	}

	supervisor := NewSupervisor()
	dispatcher := NewDispatcher(app.logger)

	var initialRunners []*TaskRunner
	for _, task := range command.Workflow {
		tr := NewTaskRunner(task, def.Env, supervisor, dispatcher)
		if len(task.When) == 0 {
			initialRunners = append(initialRunners, tr)
		} else {
			dispatcher.Register(ctx, tr)
		}
	}

	for _, tr := range initialRunners {
		tr.Run(ctx)
	}

	return supervisor.Wait()
}
