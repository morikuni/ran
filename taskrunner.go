package ran

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

type nopReceiver struct{}

func (nopReceiver) Receive(_ context.Context, _ Event) {}

var discardEvent nopReceiver

type TaskRunner struct {
	env              Env
	task             Task
	starter          WorkerStarter
	receiver         EventReceiver
	stack            Stack
	receivableTopics map[string]struct{}
	head             map[string]Event

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	logger Logger
}

func NewTaskRunner(
	task Task,
	env Env,
	starter WorkerStarter,
	receiver EventReceiver,
	stack Stack,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	logger Logger,
) *TaskRunner {
	if env == nil {
		env = Env{}
	}

	receivableTopics := make(map[string]struct{}, len(task.When))
	for _, topic := range task.When {
		receivableTopics[topic] = struct{}{}
	}

	if task.Name == "" {
		receiver = discardEvent
	}

	return &TaskRunner{
		env,
		task,
		starter,
		receiver,
		stack,
		receivableTopics,
		make(map[string]Event, len(receivableTopics)),
		stdin,
		stdout,
		stderr,
		logger,
	}
}

func (tr *TaskRunner) Run(ctx context.Context) {
	tr.run(ctx, nil)
}

func (tr *TaskRunner) run(ctx context.Context, events map[string]Event) {
	params := eventsToParams(events)
	tr.starter.Start(ctx, func(ctx context.Context) error {
		fixedEnv, err := evaluateEnv(tr.task.Env, params)
		if err != nil {
			return err
		}
		env := appendEnv(tr.env, fixedEnv)

		stdin, stdout, stderr := tr.getIO()
		if tr.task.Defer != "" {
			deferCmd := bashCmd(tr.task.Defer, stdin, stdout, stderr, env, tr.logger)
			tr.stack.Push(deferCmd)
		}

		if tr.task.Cmd == "" {
			return nil
		}

		bufOut := &bytes.Buffer{}
		bufErr := &bytes.Buffer{}
		cmd := bashCmd(tr.task.Cmd, tr.stdin, io.MultiWriter(bufOut, stdout), io.MultiWriter(bufErr, stderr), env, tr.logger)

		if err := cmd.Start(); err != nil {
			return err
		}
		pid := cmd.PID()
		tr.receiver.Receive(ctx, tr.newEvent("started", map[string]string{
			"pid": pid,
		}))

		err = cmd.Wait()
		tr.receiver.Receive(ctx, tr.newEvent("finished", map[string]string{
			"pid":    pid,
			"stdout": bufOut.String(),
			"stderr": bufErr.String(),
		}))
		if err != nil {
			tr.receiver.Receive(ctx, tr.newEvent("failed", map[string]string{
				"pid":    pid,
				"stdout": bufOut.String(),
				"stderr": bufErr.String(),
			}))
			return err
		}
		tr.receiver.Receive(ctx, tr.newEvent("succeeded", map[string]string{
			"pid":    pid,
			"stdout": bufOut.String(),
			"stderr": bufErr.String(),
		}))
		return nil
	})
}

func (tr *TaskRunner) newEvent(event string, payload map[string]string) Event {
	return NewEvent(strings.Join([]string{tr.task.Name, event}, "."), payload)
}

func (tr *TaskRunner) Receive(ctx context.Context, e Event) {
	if _, ok := tr.receivableTopics[e.Topic]; !ok {
		return
	}
	tr.head[e.Topic] = e
	if len(tr.receivableTopics) == len(tr.head) {
		events := tr.head
		tr.head = make(map[string]Event, len(tr.receivableTopics))
		tr.run(ctx, events)
	}
}

func (tr *TaskRunner) getIO() (stdin io.Reader, stdout, stderr io.Writer) {
	return tr.stdin, tr.stdout, tr.stderr
}

func eventsToParams(es map[string]Event) map[string]interface{} {
	r := make(map[string]interface{})
	for _, e := range es {
		cur := r
		paths := strings.Split(e.Topic, ".")
		for i, p := range paths {
			if i == len(paths)-1 {
				cur[p] = e.Payload
				break
			}
			m, ok := cur[p].(map[string]interface{})
			if !ok {
				m = make(map[string]interface{})
			}
			cur[p] = m
			cur = m
		}
	}
	return r
}

func evaluateEnv(env map[string]string, params map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		v, err := executeTemplate(v, params)
		if err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, nil
}

func executeTemplate(tmpl string, params map[string]interface{}) (string, error) {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func bashCmd(cmd string, stdin io.Reader, stdout, stderr io.Writer, env []string, logger Logger) Cmd {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = stdin
	c.Env = env
	return gracefulErrorCmd{cmd, c, logger}
}

type gracefulErrorCmd struct {
	cmd        string
	underlying *exec.Cmd
	logger     Logger
}

func (cmd gracefulErrorCmd) PID() string {
	return strconv.Itoa(cmd.underlying.Process.Pid)
}

func (cmd gracefulErrorCmd) log() {
	cmd.logger.Info("> %s", strings.Replace(strings.TrimRight(cmd.cmd, " \n"), "\n", "\n> ", -1))
}

func (cmd gracefulErrorCmd) Start() error {
	cmd.log()
	if err := cmd.underlying.Start(); err != nil {
		return fmt.Errorf("%q: %s", cmd.cmd, err.Error())
	}
	return nil
}

func (cmd gracefulErrorCmd) Wait() error {
	if err := cmd.underlying.Wait(); err != nil {
		return fmt.Errorf("%q: %s", cmd.cmd, err.Error())
	}
	return nil
}

func (cmd gracefulErrorCmd) Run() error {
	cmd.log()
	if err := cmd.underlying.Run(); err != nil {
		return fmt.Errorf("%q: %s", cmd.cmd, err.Error())
	}
	return nil
}
