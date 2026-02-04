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

var _ yaml.Unmarshaler = (*Root)(nil)

func (r *Root) ApplyVarsToTasks(tasks []string, vars *VarSource) (*Root, error) {
	if r == nil {
		return nil, nil
	}
	
	tLen := len(tasks)
	if tLen == 0 {
		return nil, lg.Ef("tasks len = %d")
	}
	
	globalVars := r.Var
	if vars != nil {
		globalVars = globalVars.Merge(vars, true)
	}
	
	newData := make(map[string]*Task, len(tasks))
	
	for _, name := range tasks {
		task, ok := r.Tasks.Get(name)
		if !ok {
			return nil, lg.Ef("task %s not found", name)
		}
		
		newTask, err := task.ApplyVars(globalVars)
		if err != nil {
			return nil, err
		}
		
		newData[name] = newTask
	}

	newTaskSrc := &TaskSource{}
	newTaskSrc.SetSource(newData)
	
	return &Root{
		Var: r.Var,
		Tasks: newTaskSrc,
		Env: r.Env,
	}, nil
}

func (r *Root) ApplyVarsToTask(task string, vars *VarSource) (*Root, error) {
	return r.ApplyVarsToTasks([]string{task}, vars)
}

func (r *Root) ApplyVarsAll(vars *VarSource) (*Root, error) {
	if r == nil {
		return nil, nil
	}
	
	taskNames := r.Tasks.Keys()
	
	return r.ApplyVarsToTasks(taskNames, vars)
}
