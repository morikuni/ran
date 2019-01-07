package ran

import (
	"fmt"
	"io"
	"strings"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type StdLogger struct {
	w io.Writer
}

var _ interface {
	Logger
} = (*StdLogger)(nil)

func NewStdLogger(w io.Writer) *StdLogger {
	return &StdLogger{w}
}

func (l *StdLogger) Info(msg string, args ...interface{}) {
	l.print("INFO", msg, args...)
}

func (l *StdLogger) Error(msg string, args ...interface{}) {
	l.print("ERROR", msg, args...)
}

func (l *StdLogger) print(level string, msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	msg = fmt.Sprintf("[%s] %s", level, msg)

	if len(args) == 0 {
		fmt.Fprint(l.w, msg)
		return
	}
	fmt.Fprintf(l.w, msg, args...)
}
