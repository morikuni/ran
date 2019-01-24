package ran

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

type Script interface {
	PID() string
	Start() error
	Wait() error
	Run() error
}

func shScript(script string, stdin io.Reader, stdout, stderr io.Writer, env []string, logger Logger) Script {
	c := exec.Command("sh", "-c", script)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = stdin
	c.Env = env
	return loggingScript{script, c, logger}
}

type loggingScript struct {
	script     string
	underlying *exec.Cmd
	logger     Logger
}

func (s loggingScript) PID() string {
	return strconv.Itoa(s.underlying.Process.Pid)
}

func (s loggingScript) log() {
	s.logger.Info("> %s", strings.Replace(strings.TrimRight(s.script, " \n"), "\n", "\n> ", -1))
}

func (s loggingScript) Start() error {
	s.log()
	if err := s.underlying.Start(); err != nil {
		return fmt.Errorf("%q: %s", s.script, err.Error())
	}
	return nil
}

func (s loggingScript) Wait() error {
	if err := s.underlying.Wait(); err != nil {
		return fmt.Errorf("%q: %s", s.script, err.Error())
	}
	return nil
}

func (s loggingScript) Run() error {
	s.log()
	if err := s.underlying.Run(); err != nil {
		return fmt.Errorf("%q: %s", s.script, err.Error())
	}
	return nil
}
