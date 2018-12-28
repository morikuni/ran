package workflow_test

import (
	"context"
	"testing"

	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	cases := []struct {
		task workflow.Task
		env  workflow.Env

		wantTopic  string
		wantOutput string
	}{
		{
			workflow.Task{
				Name: "simple",
				Cmd:  `echo "hello world"`,
			},
			nil,

			"simple.succeeded",
			"hello world\n",
		},
		{
			workflow.Task{
				Name: "pipe",
				Cmd:  `echo "hello world" | sed -e "s/hello/hi!/g"`,
			},
			nil,

			"pipe.succeeded",
			"hi! world\n",
		},
		{
			workflow.Task{
				Name: "command substitution backquote",
				Cmd:  "echo `echo backquote`",
			},
			nil,

			"command substitution backquote.succeeded",
			"backquote\n",
		},
		{
			workflow.Task{
				Name: "command substitution dollar",
				Cmd:  "echo $(echo dollar)",
			},
			nil,

			"command substitution dollar.succeeded",
			"dollar\n",
		},
		{
			workflow.Task{
				Name: "process substitution",
				Cmd:  "cat <(echo process)",
			},
			nil,

			"process substitution.succeeded",
			"process\n",
		},
		{
			workflow.Task{
				Name: "error",
				Cmd:  "cat nofile",
			},
			nil,

			"error.failed",
			"cat: nofile: No such file or directory\n",
		},
		{
			workflow.Task{
				Name: "env",
				Cmd:  "echo $HELLO",
			},
			workflow.Env{"HELLO=world"},

			"env.succeeded",
			"world\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.task.Name, func(t *testing.T) {
			starter := NewSynchronousStarter()
			recorder := NewEventRecorder()
			tr := workflow.NewTaskRunner(tc.task, tc.env, starter, recorder)
			tr.Run(context.Background())
			assert.NoError(t, starter.Error)
			assert.Equal(t, tc.wantTopic, recorder.GetTopic(2))
			assert.Equal(t, tc.wantOutput, recorder.GetValue(2, "output"))
		})
	}
}
