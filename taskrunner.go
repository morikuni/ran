package ran

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/template"
)

type nopReceiver struct{}

func (nopReceiver) Receive(_ Event) {}

var discardEvent nopReceiver

type TaskRunner struct {
	task          Task
	commandRunner CommandRunner

	env RuntimeEnvironment

	starter  WorkerStarter
	receiver EventReceiver
	stack    Stack

	mu               sync.Mutex
	receivableTopics map[string]struct{}
	head             map[string]Event

	logger Logger
}

func NewTaskRunner(
	task Task,
	commandRunner CommandRunner,
	starter WorkerStarter,
	receiver EventReceiver,
	stack Stack,
	logger Logger,
	env RuntimeEnvironment,
) *TaskRunner {
	receivableTopics := make(map[string]struct{}, len(task.When))
	for _, topic := range task.When {
		receivableTopics[topic] = struct{}{}
	}

	if task.Name == "" {
		receiver = discardEvent
	}

	return &TaskRunner{
		task,
		commandRunner,
		env,
		starter,
		receiver,
		stack,
		sync.Mutex{},
		receivableTopics,
		make(map[string]Event, len(receivableTopics)),
		logger,
	}
}

func (tr *TaskRunner) Run() {
	tr.run(nil)
}

func (tr *TaskRunner) run(events map[string]Event) {
	params := eventsToParams(events)
	tr.starter.Start(func() error {
		fixedEnv, err := evaluateEnv(tr.task.Env, params)
		if err != nil {
			return err
		}
		env := appendEnv(tr.env.Env, fixedEnv)

		stdin, stdout, stderr := tr.getIO()
		renv := RuntimeEnvironment{
			stdin,
			stdout,
			stderr,
			env,
			tr.env.WorkingDirectory,
		}
		tr.handleDefer(renv)
		if err := tr.handleScript(renv); err != nil {
			return err
		}
		return tr.handleCall(renv)
	})
}

func (tr *TaskRunner) handleDefer(renv RuntimeEnvironment) {
	if tr.task.Defer == "" {
		return
	}
	tr.stack.Push(bashScript(tr.task.Defer, tr.logger, renv))
}

func (tr *TaskRunner) handleScript(renv RuntimeEnvironment) error {
	if tr.task.Script == "" {
		return nil
	}
	if tr.task.Call.Command != "" {
		return fmt.Errorf("%s: cannot use both of script and call", tr.task.Script)
	}

	bufOut := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	renv.Stdout = io.MultiWriter(bufOut, renv.Stdout)
	renv.Stderr = io.MultiWriter(bufErr, renv.Stderr)
	script := bashScript(tr.task.Script, tr.logger, renv)

	if err := script.Start(); err != nil {
		return err
	}
	pid := script.PID()
	tr.receiver.Receive(tr.newEvent("started", map[string]string{
		"pid": pid,
	}))

	err := script.Wait()
	tr.receiver.Receive(tr.newEvent("finished", map[string]string{
		"pid":    pid,
		"stdout": bufOut.String(),
		"stderr": bufErr.String(),
	}))
	if err != nil {
		tr.receiver.Receive(tr.newEvent("failed", map[string]string{
			"pid":    pid,
			"stdout": bufOut.String(),
			"stderr": bufErr.String(),
		}))
		return err
	}
	tr.receiver.Receive(tr.newEvent("succeeded", map[string]string{
		"pid":    pid,
		"stdout": bufOut.String(),
		"stderr": bufErr.String(),
	}))
	return nil
}

func (tr *TaskRunner) handleCall(renv RuntimeEnvironment) error {
	if tr.task.Call.Command == "" {
		return nil
	}
	if tr.task.Script != "" {
		return fmt.Errorf("%s: cannot use both of script and call", tr.task.Call.Command)
	}

	bufOut := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	renv.Stdout = io.MultiWriter(bufOut, renv.Stdout)
	renv.Stderr = io.MultiWriter(bufErr, renv.Stderr)

	// publish started, but actually not started.
	tr.receiver.Receive(tr.newEvent("started", nil))

	err := tr.commandRunner.RunCommand(tr.task.Call.Command, renv)
	tr.receiver.Receive(tr.newEvent("finished", map[string]string{
		"stdout": bufOut.String(),
		"stderr": bufErr.String(),
	}))
	if err != nil {
		tr.receiver.Receive(tr.newEvent("failed", map[string]string{
			"stdout": bufOut.String(),
			"stderr": bufErr.String(),
		}))
		return err
	}
	tr.receiver.Receive(tr.newEvent("succeeded", map[string]string{
		"stdout": bufOut.String(),
		"stderr": bufErr.String(),
	}))
	return nil
}

func (tr *TaskRunner) newEvent(event string, payload map[string]string) Event {
	return NewEvent(strings.Join([]string{tr.task.Name, event}, "."), payload)
}

func (tr *TaskRunner) Receive(e Event) {
	if _, ok := tr.receivableTopics[e.Topic]; !ok {
		return
	}
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.head[e.Topic] = e
	if len(tr.receivableTopics) == len(tr.head) {
		events := tr.head
		tr.head = make(map[string]Event, len(tr.receivableTopics))
		tr.run(events)
	}
}

func (tr *TaskRunner) getIO() (stdin io.Reader, stdout, stderr io.Writer) {
	return tr.env.Stdin, tr.env.Stdout, tr.env.Stderr
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
