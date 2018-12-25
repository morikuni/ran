package workflow

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Definition struct {
	Env      Env
	Commands map[string]Command
}

type Command struct {
	Name     string
	Workflow []Task
}

type Task struct {
	Cmd  string
	When []string
}

type Env []string

func LoadDefinition(filename string) (Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Definition{}, err
	}
	return ParseDefinition(file)
}

func ParseDefinition(r io.Reader) (Definition, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return Definition{}, err
	}

	var raw struct {
		Env      map[string]string `yaml:"env"`
		Commands map[string]struct {
			Workflow []struct {
				Cmd  string   `yaml:"cmd"`
				When []string `yaml:"when"`
			} `yaml:"workflow"`
		} `yaml:"commands"`
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		os.Environ(),
		make(map[string]Command, len(raw.Commands)),
	}

	for name, c := range raw.Commands {
		workflow := make([]Task, len(c.Workflow))
		for i, s := range c.Workflow {
			workflow[i] = Task{
				s.Cmd,
				s.When,
			}
		}
		def.Commands[name] = Command{
			name,
			workflow,
		}
	}
	for k, v := range raw.Env {
		def.Env = append(def.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return def, nil
}
