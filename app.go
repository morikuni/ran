package ran

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App struct {
	logger Logger
}

func NewApp(logger Logger) App {
	return App{logger}
}

func (app App) Run(ctx context.Context, args []string, signal <-chan os.Signal) error {
	cmd := &cobra.Command{
		Use: "ran",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("require command")
		},
	}
	cmd.SetArgs(args[1:])

	file := cmd.PersistentFlags().StringP("file", "f", "ran.yaml", "ran definition file.")

	// parse --file flag before execute to parse and append commands.
	if err := cmd.PersistentFlags().Parse(args); err != nil && err != pflag.ErrHelp {
		cmd.Usage()
		return err
	}

	def, err := LoadDefinition(*file)
	if err != nil {
		cmd.Usage()
		return err
	}

	for _, c := range def.Commands {
		cmd.AddCommand(&cobra.Command{
			Use:   c.Name,
			Short: c.Description,
			Long:  c.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				command, ok := def.Commands[cmd.Use]
				if !ok {
					return fmt.Errorf("no such command: %s", cmd.Use)
				}

				supervisor := NewSupervisor(app.logger)
				dispatcher := NewDispatcher(app.logger)
				stack := NewStack()

				var initialRunners []*TaskRunner
				for _, task := range command.Workflow {
					tr := NewTaskRunner(task, def.Env, supervisor, dispatcher, stack, os.Stdin, os.Stdout, os.Stderr)
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
			},
		})
	}

	cmd.SilenceErrors = true
	return cmd.Execute()
}
