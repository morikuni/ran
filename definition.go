package ran

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
	Name        string
	Description string
	Workflow    []Task
}

type Task struct {
	Name string
	Cmd  string
	When []string
	Env  map[string]string
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
			Description string `yaml:"description"`
			Workflow    []struct {
				Name string            `yaml:"name"`
				Cmd  string            `yaml:"cmd"`
				When []string          `yaml:"when"`
				Env  map[string]string `yaml:"env"`
			} `yaml:"workflow"`
		} `yaml:"commands"`
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		appendEnv(os.Environ(), raw.Env),
		make(map[string]Command, len(raw.Commands)),
	}

	for name, c := range raw.Commands {
		workflow := make([]Task, len(c.Workflow))
		for i, t := range c.Workflow {
			workflow[i] = Task{
				t.Name,
				t.Cmd,
				t.When,
				t.Env,
			}
		}
		def.Commands[name] = Command{
			name,
			c.Description,
			workflow,
		}
	}

	return def, nil
}

func appendEnv(env Env, m map[string]string) Env {
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	// resize capacity to len(env) to prevent conflict when append values from multiple tasks.
	return env[:len(env):len(env)]
}
