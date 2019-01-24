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
		SilenceErrors: true,
	}
	cmd.SetArgs(args[1:])

	file := cmd.PersistentFlags().StringP("file", "f", "ran.yaml", "ran definition file.")
	logLevel := cmd.PersistentFlags().String("log-level", "error", "log level. (debug, info, error, discard)")

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

	commandRunner := NewStdCommandRunner(
		def.Commands,
		def.Env,
		os.Stdin,
		os.Stdout,
		os.Stderr,
		logger,
	)

	for _, c := range def.Commands {
		cmd.AddCommand(&cobra.Command{
			Use:   c.Name,
			Short: c.Description,
			Long:  c.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				return commandRunner.RunCommand(ctx, cmd.Use)
			},
			SilenceErrors: true,
			SilenceUsage:  true,
		})
	}

	if err := cmd.Execute(); err != nil {
		logger.Error("%s", err.Error())
		return 1
	}
	return 0
}
