package models

import "github.com/brezzgg/go-packages/lg"

type TaskSource struct {
	YamlMapSource[*Task]
}

func (t *TaskSource) ApplyVars(name string, vars *VarSource) (*TaskSource, error) {
	if t == nil {
		return nil, lg.Ef("task source is nil")
	}
	task, ok := t.Get(name)
	if !ok {
		return nil, lg.Ef("task %s not found", name)
	}

	newTask, err := task.ApplyVars(vars)
	if err != nil {
		return nil, err
	}

	newData := t.GetMapCopy()
	newData[name] = newTask

	res := &TaskSource{}
	res.SetSource(newData)

	return res, nil
}

func (t *TaskSource) ApplyVarsAll(vars *VarSource) (*TaskSource, error) {
    if t == nil {
        return nil, nil
    }
    
    newData := make(map[string]*Task, t.Len())
    
    for _, key := range t.Keys() {
        task, ok := t.Get(key)
        if !ok {
            continue
        }
        newTask, err := task.ApplyVars(vars)
        if err != nil {
            return nil, err
        }
        newData[key] = newTask
    }

	res := &TaskSource{}
	res.SetSource(newData)
    return res, nil
}

func (t *TaskSource) Clean() {
	MapClean(t.GetMap())
}

type Task struct {
	Before  *BeforeSource `yaml:"before" json:"before,omitempty"`
	After   *AfterSource  `yaml:"after" json:"after,omitempty"`
	Cmds    *CmdSource    `yaml:"cmds" json:"cmds"`
	Dir     string        `yaml:"dir" json:"dir,omitempty"`
	Vars    *VarSource    `yaml:"vars" json:"vars,omitempty"`
	Env     *EnvSource    `yaml:"envs" json:"envs,omitempty"`
	Default bool          `yaml:"default" json:"default,omitempty"`
}

func (t *Task) ApplyVars(vars *VarSource) (*Task, error) {
	if t == nil {
		return nil, nil
	}
	if vars == nil {
		vars = &VarSource{}
	}

	localVars := t.Vars
	if localVars == nil {
		localVars = &VarSource{}
	}

	merged := vars.Merge(localVars, true)
	newCmds, err := t.Cmds.ApplyVars(merged)
	if err != nil {
		return nil, err
	}

	return &Task{
		Before:  t.Before,
		After:   t.After,
		Cmds:    newCmds,
		Dir:     t.Dir,
		Vars:    t.Vars,
		Env:     t.Env,
		Default: t.Default,
	}, nil
}

type AfterSource struct {
	YamlSliceSource[string]
}

type BeforeSource struct {
	YamlSliceSource[string]
}
