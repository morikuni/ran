package ran_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/morikuni/ran"
	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	cases := map[string]struct {
		task ran.Task
		env  ran.Env

		wantTopics []string
		wantStdout string
		wantStderr string
	}{
		"success": {
			ran.Task{
				Name:   "success",
				Script: `echo "$VALUE"`,
				Env: map[string]string{
					"VALUE": "hello world",
				},
			},
			nil,

			[]string{"success.started", "success.finished", "success.succeeded"},
			"hello world\n",
			"",
		},
		"error": {
			ran.Task{
				Name:   "error",
				Script: "cat nofile",
			},
			nil,

			[]string{"error.started", "error.finished", "error.failed"},
			"",
			"cat: nofile: No such file or directory\n",
		},
		"defer": {
			ran.Task{
				Name:  "defer",
				Defer: "echo defer",
			},
			nil,

			nil,
			"defer\n",
			"",
		},
		"no events": {
			ran.Task{
				Script: "echo no name",
			},
			nil,

			nil,
			"no name\n",
			"",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			starter := NewSynchronousStarter()
			recorder := NewEventRecorder()
			stack := ran.NewStack()
			stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
			logger := ran.NewStdLogger(ioutil.Discard, ran.Discard)
			tr := ran.NewTaskRunner(tc.task, tc.env, starter, recorder, stack, bytes.NewReader(nil), stdout, stderr, logger)
			tr.Run(context.Background())
			for {
				script, ok := stack.Pop()
				if !ok {
					break
				}
				require.NoError(t, script.Run())
			}
			var topics []string
			for _, e := range recorder.Events {
				topics = append(topics, e.Topic)
			}

			assert.NoError(t, starter.Error)
			assert.Equal(t, tc.wantTopics, topics)
			assert.Equal(t, tc.wantStdout, stdout.String())
			assert.Equal(t, tc.wantStderr, stderr.String())
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
