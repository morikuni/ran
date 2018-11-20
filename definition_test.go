package workflow_test

import (
	"github.com/morikuni/workflow"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseDefinition(t *testing.T) {
	r := strings.NewReader(`
tasks:
- name: aaa
  cmd: bbb
- name: ccc
  cmd: ddd
`)
	want := workflow.Definition{
		Tasks: []workflow.Task{
			{
				Name: "aaa",
				CMD: "bbb",
			},
			{
				Name: "ccc",
				CMD: "ddd",
			},
		},
	}

	def, err := workflow.ParseDefinition(r)
	assert.NoError(t, err)
	assert.Equal(t, want, def)
}

func TestLoadDefinition(t *testing.T) {
	want := workflow.Definition{
		Tasks: []workflow.Task{
			{
				Name: "echo task",
				CMD: `echo "hello"`,
			},
			{
				Name: "pipe task",
				CMD: `echo "world" | cat`,
			},
		},
	}

	def, err := workflow.LoadDefinition("testdata/simple.yaml")
	assert.NoError(t, err)
	assert.Equal(t, want, def)
}
