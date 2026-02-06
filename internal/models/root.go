package models

import (
	"github.com/brezzgg/go-packages/lg"
	"gopkg.in/yaml.v3"
)

type Root struct {
	Var   *VarSource  `yaml:"vars" json:"vars,omitempty"`
	Tasks *TaskSource `yaml:"tasks" json:"tasks"`
	Env   *EnvSource  `yaml:"envs" json:"envs,omitempty"`
}

func (r *Root) UnmarshalYAML(value *yaml.Node) error {
	type rootAlias Root
	var root rootAlias

	if err := value.Decode(&root); err != nil {
		return err
	}

	root.Tasks.Clean()
	if root.Var == nil {
		root.Var = &VarSource{}
	}
	if root.Env == nil {
		root.Env = &EnvSource{}
	}

	r.Tasks = root.Tasks
	r.Var = root.Var
	r.Env = root.Env

	return nil
}

func (r *Root) ApplyVars(vars *VarSource) (*Root, error) {
	if r == nil {
		return nil, lg.Ef("root is nil")
	}
	if r.Tasks == nil {
		return &Root{
			Var:     r.Var,
			Tasks:   r.Tasks,
			Env:     r.Env,
			applied: true,
		}, nil
	}
	global := r.Var
	if vars != nil {
		global = global.Merge(vars, true)
	}
	tasks := make(map[string]*Task, r.Tasks.Len())
	for k, v := range r.Tasks.GetMapCopy() {
		applied, err := v.ApplyVars(global)
		if err != nil {
			return nil, lg.Ef("task %s: %w", k, err)
		}
		tasks[k] = applied
	}
	taskSrc := &TaskSource{}
	taskSrc.SetSource(tasks)
	return &Root{
		Var:     r.Var,
		Tasks:   taskSrc,
		Env:     r.Env,
		applied: true,
	}, nil
}

func (r *Root) Applied() bool {
	return r.applied
}

	}
	
	taskNames := r.Tasks.Keys()
	
	return r.ApplyVarsToTasks(taskNames, vars)
}
