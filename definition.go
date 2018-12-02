package workflow

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Definition struct {
	Env      Env
	Tasks    map[string]Task
	Commands map[string]Command
}

type Task struct {
	Name string
	CMD  string
}

type Command struct {
	Name     string
	Workflow []Work
}

type Work struct {
	Run string
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
		Env   map[string]string `yaml:"env"`
		Tasks map[string]struct {
			CMD string `yaml:"cmd"`
		} `yaml:"tasks"`
		Commands map[string]struct {
			Workflow []struct {
				Run string `yaml:"run"`
			} `yaml:"workflow"`
		} `yaml:"commands"`
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		os.Environ(),
		make(map[string]Task, len(raw.Tasks)),
		make(map[string]Command, len(raw.Commands)),
	}

	for name, t := range raw.Tasks {
		def.Tasks[name] = Task{
			name,
			t.CMD,
		}
	}
	for name, c := range raw.Commands {
		works := make([]Work, len(c.Workflow))
		for i, s := range c.Workflow {
			works[i] = Work{
				s.Run,
			}
		}
		def.Commands[name] = Command{
			name,
			works,
		}
	}
	for k, v := range raw.Env {
		def.Env = append(def.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return def, nil
}
