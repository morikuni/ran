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
	app := ran.NewApp()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)

	return app.Run(context.Background(), os.Args, sig)
}
