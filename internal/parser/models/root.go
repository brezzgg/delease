package models

import (
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

type Root struct {
	Vars  Vars  `yaml:"vars"`
	Tasks Tasks `yaml:"tasks"`
}

func (r *Root) GetVars() Vars {
	return r.Vars
}

func (r *Root) GetTask(name string) (*Task, error) {
	if r.Tasks == nil {
		return nil, lg.Ef("tasks is empty")
	}

	task, ok := r.Tasks[name]
	if !ok {
		return nil, lg.Ef("task '%s' not found", name)
	}

	if task == nil {
		return nil, lg.Ef("task '%s' is nil", name)
	}

	return task, nil
}

func (r *Root) UnmarshalYAML(node *yaml.Node) error {
	type rootAlias Root

	var alias rootAlias
	if err := node.Decode(&alias); err != nil {
		return err
	}

	r.Vars = alias.Vars
	r.Tasks = alias.Tasks

	if r.Tasks != nil {
		cleanMap(r.Tasks)
	}

	return nil
}

type (
	Vars  map[string]string
	Tasks map[string]*Task
)

type Variabler interface {
	GetVars() Vars
}
