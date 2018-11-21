package main

import (
	"fmt"
	"log"
	"os"

	"github.com/morikuni/workflow"
)

func main() {
	os.Exit(app())
}

func app() (exitCode int) {
	def, err := workflow.LoadDefinition(os.Args[1])
	if err != nil {
		log.Println(err)
		return 1
	}

	for _, task := range def.Tasks {
		fmt.Println("[" + task.Name + "]")
		tr := workflow.NewTaskRunner()
		if err := tr.Run(task); err != nil {
			log.Println(err)
			return 1
		}
		fmt.Print(tr.Output())
	}

	return
}
