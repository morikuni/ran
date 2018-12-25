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

vars:
  cmd: &bbb bbb

commands:
  test:
    workflow:
    - cmd: aaa
    - cmd: *bbb
`)
	want := workflow.Definition{
		Commands: map[string]workflow.Command{
			"test": {
				Name: "test",
				Workflow: []workflow.Task{
					{
						Cmd: "aaa",
					},
					{
						Cmd: "bbb",
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
		Commands: map[string]workflow.Command{
			"all": {
				Name: "all",
				Workflow: []workflow.Task{
					{
						Cmd: `echo "hello"`,
					},
					{
						Cmd: `echo "world" | cat`,
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
