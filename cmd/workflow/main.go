package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/morikuni/workflow"
)

func main() {
	os.Exit(app())
}

func app() (exitCode int) {
	file := flag.String("f", "workflow.yaml", "file")
	flag.Parse()

	def, err := workflow.LoadDefinition(*file)
	if err != nil {
		log.Println(err)
		return 1
	}

	target := flag.Arg(0)
	command, ok := def.Commands[target]
	if !ok {
		log.Println("no such workflow:", target)
		return 1
	}

	for _, work := range command.Workflow {
		fmt.Println("[" + work.Run + "]")

		task, ok := def.Tasks[work.Run]
		if !ok {
			log.Println("no such task:", work.Run)
			return 1
		}

		tr := workflow.NewTaskRunner(def.Env)
		if err := tr.Run(task); err != nil {
			log.Println(err)
			fmt.Print(tr.Output())
			return 1
		}
		fmt.Print(tr.Output())
	}

	return
}
