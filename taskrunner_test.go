package workflow_test

import (
	"testing"

	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	cases := map[string]struct {
		task workflow.Task

		wantOutput string
		wantErr    bool
	}{
		"simple": {
			workflow.Task{
				CMD: `echo "hello world"`,
			},

			"hello world\n",
			false,
		},
		"pipe": {
			workflow.Task{
				CMD: `echo "hello world" | sed -e "s/hello/hi!/g"`,
			},

			"hi! world\n",
			false,
		},
		"command substitution backquote": {
			workflow.Task{
				CMD: "echo `echo backquote`",
			},

			"backquote\n",
			false,
		},
		"command substitution dollar": {
			workflow.Task{
				CMD: "echo $(echo dollar)",
			},

			"dollar\n",
			false,
		},
		"process substitution": {
			workflow.Task{
				CMD: "cat <(echo process)",
			},

			"process\n",
			false,
		},
		"error": {
			workflow.Task{
				CMD: "cat nofile",
			},

			"cat: nofile: No such file or directory\n",
			true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tr := workflow.NewTaskRunner()
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
