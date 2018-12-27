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

		wantOutput string
		wantErr    bool
	}{
		{
			workflow.Task{
				Name: "simple",
				Cmd:  `echo "hello world"`,
			},
			nil,

			"hello world\n",
			false,
		},
		{
			workflow.Task{
				Name: "pipe",
				Cmd:  `echo "hello world" | sed -e "s/hello/hi!/g"`,
			},
			nil,

			"hi! world\n",
			false,
		},
		{
			workflow.Task{
				Name: "command substitution backquote",
				Cmd:  "echo `echo backquote`",
			},
			nil,

			"backquote\n",
			false,
		},
		{
			workflow.Task{
				Name: "command substitution dollar",
				Cmd:  "echo $(echo dollar)",
			},
			nil,

			"dollar\n",
			false,
		},
		{
			workflow.Task{
				Name: "process substitution",
				Cmd:  "cat <(echo process)",
			},
			nil,

			"process\n",
			false,
		},
		{
			workflow.Task{
				Name: "error",
				Cmd:  "cat nofile",
			},
			nil,

			"cat: nofile: No such file or directory\n",
			true,
		},
		{
			workflow.Task{
				Name: "env",
				Cmd:  "echo $HELLO",
			},
			workflow.Env{"HELLO=world"},

			"world\n",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.task.Name, func(t *testing.T) {
			starter := NewSynchronousStarter()
			recorder := NewEventRecorder()
			tr := workflow.NewTaskRunner(tc.task, tc.env, starter, recorder)
			tr.Run(context.Background())
			if tc.wantErr {
				assert.Equal(t, tc.task.Name+".failed", recorder.GetTopic(0))
			} else {
				assert.Equal(t, tc.task.Name+".succeeded", recorder.GetTopic(0))
			}
			assert.Equal(t, tc.wantOutput, recorder.GetValue(0, "output"))
		})
	}
}
