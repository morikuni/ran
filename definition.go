package workflow

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Definition struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name string `yaml:"name"`
	CMD  string `yaml:"cmd"`
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
	var def Definition
	if err := yaml.Unmarshal(bs, &def); err != nil {
		return Definition{}, err
	}
	return def, nil
}
