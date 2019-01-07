package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/morikuni/ran"
)

func main() {
	os.Exit(app())
}

func app() (exitCode int) {
	logger := ran.NewStdLogger(os.Stdout)
	app := ran.NewApp(logger)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)

	err := app.Run(context.Background(), os.Args, sig)
	if err != nil {
		logger.Error(err.Error())
		return 1
	}

	return 0
}
