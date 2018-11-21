package workflow

import (
	"bytes"
	"os/exec"
)

type TaskRunner struct {
	output *bytes.Buffer
}

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{
		&bytes.Buffer{},
	}
}

func (tr *TaskRunner) Run(task Task) error {
	cmd := exec.Command("bash", "-c", task.CMD)
	cmd.Stdout = tr.output
	cmd.Stderr = tr.output
	return cmd.Run()
}

func (tr *TaskRunner) Output() string {
	return tr.output.String()
}
