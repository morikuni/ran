package workflow

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Definition struct {
	Tasks map[string]Task
}

type Task struct {
	Name string
	CMD  string
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
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		make(map[string]Task, len(raw.Tasks)),
	}
	for name, t := range raw.Tasks {
		def.Tasks[name] = Task{
			name,
			t.CMD,
		}
	}

	return def, nil
}
