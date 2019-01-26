package main

import (
	"os"

	"github.com/morikuni/ran"
)

func main() {
	app := ran.NewApp()
	os.Exit(app.Run(os.Args, os.Stdin, os.Stdout, os.Stderr))
}
