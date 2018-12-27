package workflow

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type TaskRunner struct {
	env      Env
	task     Task
	starter  WorkerStarter
	receiver EventReceiver
}

func NewTaskRunner(task Task, env Env, starter WorkerStarter, receiver EventReceiver) *TaskRunner {
	if env == nil {
		env = Env{}
	}
	return &TaskRunner{
		env,
		task,
		starter,
		receiver,
	}
}

func (tr *TaskRunner) Run(ctx context.Context) {
	tr.starter.Start(ctx, func(ctx context.Context) error {
		cmd := exec.Command("bash", "-c", tr.task.Cmd)

		buf := &bytes.Buffer{}
		cmd.Stdout = io.MultiWriter(buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(buf, os.Stderr)
		cmd.Env = tr.env

		err := cmd.Run()
		if err != nil {
			return tr.receiver.Receive(ctx, tr.newEvent("failed", map[string]string{
				"output": buf.String(),
			}))
		}
		return tr.receiver.Receive(ctx, tr.newEvent("succeeded", map[string]string{
			"output": buf.String(),
		}))
	})
}

func (tr *TaskRunner) newEvent(event string, payload map[string]string) Event {
	return NewEvent(strings.Join([]string{tr.task.Name, event}, "."), payload)
}

func (tr *TaskRunner) Receive(ctx context.Context, e Event) error {
	return nil
}
