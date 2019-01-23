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
}

func NewApp() App {
	return App{}
}

func (app App) Run(ctx context.Context, args []string, signal <-chan os.Signal) int {
	cmd := &cobra.Command{
		Use: "ran",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("require command")
		},
	}
	cmd.SetArgs(args[1:])

	file := cmd.PersistentFlags().StringP("file", "f", "ran.yaml", "ran definition file.")
	logLevel := cmd.PersistentFlags().String("log-level", "discard", "log level. (debug, info, error, discard)")

	// parse --file flag before execute to parse and append commands.
	if err := cmd.PersistentFlags().Parse(args); err != nil && err != pflag.ErrHelp {
		cmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	level, err := NewLogLevel(*logLevel)
	if err != nil {
		cmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	logger := NewStdLogger(os.Stdout, level)

	def, err := LoadDefinition(*file)
	if err != nil {
		cmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
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

				supervisor := NewSupervisor(logger)
				dispatcher := NewDispatcher(logger)
				stack := NewStack()

				var initialRunners []*TaskRunner
				for _, task := range command.Tasks {
					tr := NewTaskRunner(task, def.Env, supervisor, dispatcher, stack, os.Stdin, os.Stdout, os.Stderr, logger)
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
	if err := cmd.Execute(); err != nil {
		logger.Error("%s", err.Error())
		return 1
	}
	return 0
}
