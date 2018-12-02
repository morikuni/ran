package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/morikuni/workflow"
)

func main() {
	os.Exit(app())
}

func app() (exitCode int) {
	logger := workflow.NewStdLogger(os.Stdout)
	app := workflow.NewApp(logger)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)

	err := app.Run(context.Background(), os.Args, sig)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
