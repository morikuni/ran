package workflow_test

import (
	"testing"

	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	cases := map[string]struct {
		env  workflow.Env
		task workflow.Task

		wantOutput string
		wantErr    bool
	}{
		"simple": {
			nil,
			workflow.Task{
				Cmd: `echo "hello world"`,
			},

			"hello world\n",
			false,
		},
		"pipe": {
			nil,
			workflow.Task{
				Cmd: `echo "hello world" | sed -e "s/hello/hi!/g"`,
			},

			"hi! world\n",
			false,
		},
		"command substitution backquote": {
			nil,
			workflow.Task{
				Cmd: "echo `echo backquote`",
			},

			"backquote\n",
			false,
		},
		"command substitution dollar": {
			nil,
			workflow.Task{
				Cmd: "echo $(echo dollar)",
			},

			"dollar\n",
			false,
		},
		"process substitution": {
			nil,
			workflow.Task{
				Cmd: "cat <(echo process)",
			},

			"process\n",
			false,
		},
		"error": {
			nil,
			workflow.Task{
				Cmd: "cat nofile",
			},

			"cat: nofile: No such file or directory\n",
			true,
		},
		"env": {
			workflow.Env{"HELLO=world"},
			workflow.Task{
				Cmd: "echo $HELLO",
			},

			"world\n",
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tr := workflow.NewTaskRunner(tc.env)
			err := tr.Run(tc.task)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantOutput, tr.Output())
		})
	}
}
