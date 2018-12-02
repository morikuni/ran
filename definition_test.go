package workflow_test

import (
	"strings"
	"testing"

	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
)

func TestParseDefinition(t *testing.T) {
	r := strings.NewReader(`
env:
  FOO: 123
  HELLO: world

tasks:
  aaa:
    cmd: bbb
  ccc:
    cmd: ddd

commands:
  test:
    workflow:
    - run: aaa
    - run: ccc
`)
	want := workflow.Definition{
		Tasks: map[string]workflow.Task{
			"aaa": {
				Name: "aaa",
				CMD:  "bbb",
			},
			"ccc": {
				Name: "ccc",
				CMD:  "ddd",
			},
		},
		Commands: map[string]workflow.Command{
			"test": {
				Name: "test",
				Workflow: []workflow.Work{
					{
						Run: "aaa",
					},
					{
						Run: "ccc",
					},
				},
			},
		},
	}

	def, err := workflow.ParseDefinition(r)
	env := def.Env
	def.Env = nil
	assert.NoError(t, err)
	assert.Equal(t, want, def)
	assert.Contains(t, env, "FOO=123")
	assert.Contains(t, env, "HELLO=world")
}

func TestLoadDefinition(t *testing.T) {
	want := workflow.Definition{
		Tasks: map[string]workflow.Task{
			"echo": {
				Name: "echo",
				CMD:  `echo "hello"`,
			},
			"pipe": {
				Name: "pipe",
				CMD:  `echo "world" | cat`,
			},
		},
		Commands: map[string]workflow.Command{
			"all": {
				Name: "all",
				Workflow: []workflow.Work{
					{
						Run: "echo",
					},
					{
						Run: "pipe",
					},
				},
			},
		},
	}

	def, err := workflow.LoadDefinition("testdata/simple.yaml")
	env := def.Env
	def.Env = nil
	assert.NoError(t, err)
	assert.Equal(t, want, def)
	assert.Contains(t, env, "ENV_X=123")
	assert.Contains(t, env, "ENV_Y=hello")
}
