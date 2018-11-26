package workflow

import (
	"bytes"
	"os/exec"
)

type TaskRunner struct {
	env    Env
	output *bytes.Buffer
}

func NewTaskRunner(env Env) *TaskRunner {
	if env == nil {
		env = Env{}
	}
	return &TaskRunner{
		env,
		&bytes.Buffer{},
	}
}

func (tr *TaskRunner) Run(task Task) error {
	cmd := exec.Command("bash", "-c", task.CMD)
	cmd.Env = tr.env
	cmd.Stdout = tr.output
	cmd.Stderr = tr.output
	return cmd.Run()
}

func (tr *TaskRunner) Output() string {
	return tr.output.String()
}
