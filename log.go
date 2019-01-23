package ran

import (
	"fmt"
	"io"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func NewLogLevel(level string) (LogLevel, error) {
	l, ok := map[string]LogLevel{
		"debug":   Debug,
		"info":    Info,
		"error":   Error,
		"discard": Discard,
	}[level]
	if !ok {
		return 0, fmt.Errorf("unknown log level: %q", level)
	}
	return l, nil
}

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
	Discard
)

func (l LogLevel) String() string {
	return map[LogLevel]string{
		Debug:   "DEBUG",
		Info:    "INFO",
		Error:   "ERROR",
		Discard: "DISCARD",
	}[l]
}

type StdLogger struct {
	w     io.Writer
	level LogLevel
}

var _ interface {
	Logger
} = (*StdLogger)(nil)

func NewStdLogger(w io.Writer, level LogLevel) *StdLogger {
	return &StdLogger{w, level}
}

func (l *StdLogger) Debug(msg string, args ...interface{}) {
	l.print(Debug, msg, args...)
}

func (l *StdLogger) Info(msg string, args ...interface{}) {
	l.print(Info, msg, args...)
}

func (l *StdLogger) Error(msg string, args ...interface{}) {
	l.print(Error, msg, args...)
}

func (l *StdLogger) print(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	msg = fmt.Sprintf("%s", msg)

	if len(args) == 0 {
		fmt.Fprint(l.w, msg)
		return
	}
	fmt.Fprintf(l.w, msg, args...)
}
