package workflow

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Definition struct {
	Tasks    map[string]Task
	Workflow map[string][]Stage
}

type Task struct {
	Name string
	CMD  string
}

type Stage struct {
	Run string
}

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
		Tasks map[string]struct {
			CMD string `yaml:"cmd"`
		} `yaml:"tasks"`
		Workflow map[string][]struct {
			Run string `yaml:"run"`
		} `yaml:"workflow"`
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		make(map[string]Task, len(raw.Tasks)),
		make(map[string][]Stage, len(raw.Workflow)),
	}
	for name, t := range raw.Tasks {
		def.Tasks[name] = Task{
			name,
			t.CMD,
		}
	}
	for name, ss := range raw.Workflow {
		stages := make([]Stage, len(ss))
		for i, s := range ss {
			stages[i] = Stage{
				s.Run,
			}
		}
		def.Workflow[name] = stages
	}

	return def, nil
}
