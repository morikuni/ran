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
	Tasks       []Task
}

type Task struct {
	Name   string
	Script string
	When   []string
	Env    map[string]string
	Defer  string
}

type Env []string

func LoadDefinition(filename string) (Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Definition{}, err
	}
	defer file.Close()
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
			Tasks       []struct {
				Name   string            `yaml:"name"`
				Script string            `yaml:"script"`
				When   []string          `yaml:"when"`
				Env    map[string]string `yaml:"env"`
				Defer  string            `yaml:"defer"`
			} `yaml:"tasks"`
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
		tasks := make([]Task, len(c.Tasks))
		for i, t := range c.Tasks {
			tasks[i] = Task{
				t.Name,
				t.Script,
				t.When,
				t.Env,
				t.Defer,
			}
		}
		def.Commands[name] = Command{
			name,
			c.Description,
			tasks,
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
