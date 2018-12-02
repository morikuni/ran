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

	for _, work := range command.Workflow {
		app.logger.Info("[" + work.Run + "]")

		task, ok := def.Tasks[work.Run]
		if !ok {
			return fmt.Errorf("no such task: %s", work.Run)
		}

		tr := NewTaskRunner(def.Env)
		if err := tr.Run(task); err != nil {
			app.logger.Info(tr.Output())
			return err
		}
		app.logger.Info(tr.Output())
	}

	return nil
}
