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

func bashScript(script string, logger Logger, env RuntimeEnvironment) Script {
	c := exec.Command("bash", "-c", script)
	c.Stdin = env.Stdin
	c.Stdout = env.Stdout
	c.Stderr = env.Stderr
	c.Env = env.Env
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

type RuntimeEnvironment struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    EnvironmentVariables
}
