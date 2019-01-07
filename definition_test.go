package ran_test

import (
	"strings"
	"testing"

	"github.com/morikuni/ran"
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
	want := ran.Definition{
		Commands: map[string]ran.Command{
			"test": {
				Name: "test",
				Workflow: []ran.Task{
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

	def, err := ran.ParseDefinition(r)
	env := def.Env
	def.Env = nil
	assert.NoError(t, err)
	assert.Equal(t, want, def)
	assert.Contains(t, env, "FOO=123")
	assert.Contains(t, env, "HELLO=world")
}

func TestLoadDefinition(t *testing.T) {
	want := ran.Definition{
		Commands: map[string]ran.Command{
			"all": {
				Name: "all",
				Workflow: []ran.Task{
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

	def, err := ran.LoadDefinition("testdata/simple.yaml")
	env := def.Env
	def.Env = nil
	assert.NoError(t, err)
	assert.Equal(t, want, def)
	assert.Contains(t, env, "ENV_X=123")
	assert.Contains(t, env, "ENV_Y=hello")
}
