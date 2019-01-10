package ran_test

import (
	"context"
	"testing"

	"github.com/morikuni/ran"
	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	cases := []struct {
		task ran.Task
		env  ran.Env

		wantTopic  string
		wantStdout string
		wantStderr string
	}{
		{
			ran.Task{
				Name: "simple",
				Cmd:  `echo "hello world"`,
			},
			nil,

			"simple.succeeded",
			"hello world\n",
			"",
		},
		{
			ran.Task{
				Name: "pipe",
				Cmd:  `echo "hello world" | sed -e "s/hello/hi!/g"`,
			},
			nil,

			"pipe.succeeded",
			"hi! world\n",
			"",
		},
		{
			ran.Task{
				Name: "command substitution backquote",
				Cmd:  "echo `echo backquote`",
			},
			nil,

			"command substitution backquote.succeeded",
			"backquote\n",
			"",
		},
		{
			ran.Task{
				Name: "command substitution dollar",
				Cmd:  "echo $(echo dollar)",
			},
			nil,

			"command substitution dollar.succeeded",
			"dollar\n",
			"",
		},
		{
			ran.Task{
				Name: "process substitution",
				Cmd:  "cat <(echo process)",
			},
			nil,

			"process substitution.succeeded",
			"process\n",
			"",
		},
		{
			ran.Task{
				Name: "error",
				Cmd:  "cat nofile",
			},
			nil,

			"error.failed",
			"",
			"cat: nofile: No such file or directory\n",
		},
		{
			ran.Task{
				Name: "env",
				Cmd:  "echo $HELLO",
			},
			ran.Env{"HELLO=world"},

			"env.succeeded",
			"world\n",
			"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.task.Name, func(t *testing.T) {
			starter := NewSynchronousStarter()
			recorder := NewEventRecorder()
			tr := ran.NewTaskRunner(tc.task, tc.env, starter, recorder)
			tr.Run(context.Background())
			assert.NoError(t, starter.Error)
			assert.Equal(t, tc.wantTopic, recorder.GetTopic(2))
			assert.Equal(t, tc.wantStdout, recorder.GetValue(2, "stdout"))
			assert.Equal(t, tc.wantStderr, recorder.GetValue(2, "stderr"))
		})
	}
}

func Test_EventsToParams(t *testing.T) {
	events := map[string]ran.Event{
		"aa.bb": ran.Event{
			Topic:   "aa.bb",
			Payload: map[string]string{"value1": "1"},
		},
		"aa.cc": ran.Event{
			Topic:   "aa.cc",
			Payload: map[string]string{"value2": "2"},
		},
		"xx.yy": ran.Event{
			Topic:   "xx.yy",
			Payload: map[string]string{"value3": "3"},
		},
	}

	m := ran.EventsToParams(events)
	expect := map[string]interface{}{
		"aa": map[string]interface{}{
			"bb": map[string]string{
				"value1": "1",
			},
			"cc": map[string]string{
				"value2": "2",
			},
		},
		"xx": map[string]interface{}{
			"yy": map[string]string{
				"value3": "3",
			},
		},
	}
	assert.Equal(t, expect, m)
}
